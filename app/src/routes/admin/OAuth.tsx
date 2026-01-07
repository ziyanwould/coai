import { useState, useEffect } from "react";
import { useTranslation } from "react-i18next";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Switch } from "@/components/ui/switch";
import { Label } from "@/components/ui/label";
import { useToast } from "@/components/ui/use-toast";
import { toastState } from "@/api/common.ts";
import {
  getOAuthConfig,
  updateOAuthConfig,
  LinuxDoOAuth,
} from "@/admin/api/oauth.ts";
import { Loader2, Save, Info, KeyRound, ExternalLink } from "lucide-react";
import { Alert, AlertDescription } from "@/components/ui/alert";

function OAuth() {
  const { t } = useTranslation();
  const { toast } = useToast();
  const [loading, setLoading] = useState<boolean>(true);
  const [saving, setSaving] = useState<boolean>(false);

  // Linux.do OAuth 配置
  const [linuxDoEnabled, setLinuxDoEnabled] = useState<boolean>(false);
  const [linuxDoClientId, setLinuxDoClientId] = useState<string>("");
  const [linuxDoClientSecret, setLinuxDoClientSecret] = useState<string>("");
  const [linuxDoRedirectUrl, setLinuxDoRedirectUrl] = useState<string>("");
  const [hasExistingSecret, setHasExistingSecret] = useState<boolean>(false);

  // 加载配置
  useEffect(() => {
    async function loadConfig() {
      setLoading(true);
      try {
        const res = await getOAuthConfig();
        if (res.status && res.data) {
          const linuxDo = res.data.linux_do;
          setLinuxDoEnabled(linuxDo.enabled || false);
          setLinuxDoClientId(linuxDo.client_id || "");
          setLinuxDoRedirectUrl(linuxDo.redirect_url || "");
          setHasExistingSecret(linuxDo.has_secret || false);
        } else if (res.error) {
          toast({
            title: t("admin.oauth.load-error", "加载 OAuth 配置失败"),
            description: res.error,
            variant: "destructive",
          });
        }
      } catch (error) {
        console.error("Failed to load OAuth config:", error);
        toast({
          title: t("admin.oauth.load-error", "加载 OAuth 配置失败"),
          description: String(error),
          variant: "destructive",
        });
      } finally {
        setLoading(false);
      }
    }
    loadConfig();
  }, [t, toast]);

  // 保存配置
  const handleSave = async () => {
    setSaving(true);

    const linuxDoConfig: LinuxDoOAuth = {
      enabled: linuxDoEnabled,
      client_id: linuxDoClientId,
      client_secret: linuxDoClientSecret, // 空字符串时后端会保留旧值
      redirect_url: linuxDoRedirectUrl,
    };

    const res = await updateOAuthConfig(linuxDoConfig);
    toastState(toast, t, res, true);

    if (res.status) {
      // 保存成功后更新状态
      if (linuxDoClientSecret) {
        setHasExistingSecret(true);
        setLinuxDoClientSecret(""); // 清空输入的密钥
      }
    }

    setSaving(false);
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Loader2 className="w-8 h-8 animate-spin" />
      </div>
    );
  }

  return (
    <div className="oauth-page p-4 max-w-4xl mx-auto">
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <KeyRound className="w-5 h-5" />
            {t("admin.oauth.title", "OAuth 配置")}
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-6">
          <Alert>
            <Info className="h-4 w-4" />
            <AlertDescription>
              {t(
                "admin.oauth.description",
                "配置第三方 OAuth 登录，允许用户使用外部账号登录系统。"
              )}
            </AlertDescription>
          </Alert>

          {/* Linux.do OAuth 配置 */}
          <div className="border rounded-lg p-4 space-y-4">
            <div className="flex items-center justify-between">
              <div className="space-y-0.5">
                <div className="flex items-center gap-2">
                  <Label
                    htmlFor="linux-do-enabled"
                    className="text-base font-medium"
                  >
                    Linux.do OAuth
                  </Label>
                  <a
                    href="https://connect.linux.do"
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-muted-foreground hover:text-primary"
                  >
                    <ExternalLink className="w-4 h-4" />
                  </a>
                </div>
                <p className="text-sm text-muted-foreground">
                  {t(
                    "admin.oauth.linux-do-desc",
                    "启用 Linux.do 社区账号登录"
                  )}
                </p>
              </div>
              <Switch
                id="linux-do-enabled"
                checked={linuxDoEnabled}
                onCheckedChange={setLinuxDoEnabled}
              />
            </div>

            <div className="space-y-3">
              <div className="space-y-2">
                <Label htmlFor="client-id">Client ID</Label>
                <Input
                  id="client-id"
                  value={linuxDoClientId}
                  onChange={(e) => setLinuxDoClientId(e.target.value)}
                  placeholder={t(
                    "admin.oauth.client-id-placeholder",
                    "输入 Client ID"
                  )}
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="client-secret">
                  Client Secret
                  {hasExistingSecret && (
                    <span className="ml-2 text-xs text-green-600">
                      {t("admin.oauth.secret-configured", "(已配置)")}
                    </span>
                  )}
                </Label>
                <Input
                  id="client-secret"
                  type="password"
                  value={linuxDoClientSecret}
                  onChange={(e) => setLinuxDoClientSecret(e.target.value)}
                  placeholder={
                    hasExistingSecret
                      ? t(
                          "admin.oauth.secret-placeholder-update",
                          "留空保持不变，输入新值则更新"
                        )
                      : t(
                          "admin.oauth.secret-placeholder",
                          "输入 Client Secret"
                        )
                  }
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="redirect-url">Redirect URL</Label>
                <Input
                  id="redirect-url"
                  value={linuxDoRedirectUrl}
                  onChange={(e) => setLinuxDoRedirectUrl(e.target.value)}
                  placeholder={t(
                    "admin.oauth.redirect-placeholder",
                    "例如: https://your-domain.com/oauth-success"
                  )}
                />
                <p className="text-xs text-muted-foreground">
                  {t(
                    "admin.oauth.redirect-hint",
                    "请确保此 URL 与 OAuth 应用配置中的回调地址一致"
                  )}
                </p>
              </div>
            </div>
          </div>

          <div className="flex gap-3 pt-4">
            <Button
              onClick={handleSave}
              disabled={saving}
              className="flex items-center gap-2"
            >
              {saving ? (
                <>
                  <Loader2 className="w-4 h-4 animate-spin" />
                  {t("admin.oauth.saving", "保存中...")}
                </>
              ) : (
                <>
                  <Save className="w-4 h-4" />
                  {t("admin.oauth.save", "保存配置")}
                </>
              )}
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}

export default OAuth;
