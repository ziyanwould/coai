import { useState, useEffect } from "react";
import { useTranslation } from "react-i18next";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { useToast } from "@/components/ui/use-toast";
import { toastState } from "@/api/common.ts";
import {
  getVisionConfig,
  updateVisionConfig,
  refreshVisionConfig,
} from "@/admin/api/vision.ts";
import { Loader2, Save, RefreshCw, Info } from "lucide-react";
import { Alert, AlertDescription } from "@/components/ui/alert";

function Vision() {
  const { t } = useTranslation();
  const { toast } = useToast();
  const [models, setModels] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(true);
  const [saving, setSaving] = useState<boolean>(false);
  const [refreshing, setRefreshing] = useState<boolean>(false);

  // 加载配置
  useEffect(() => {
    async function loadConfig() {
      setLoading(true);
      try {
        const res = await getVisionConfig();
        if (res.status && res.data) {
          setModels(res.data.models ? res.data.models.join("\n") : "");
        } else if (res.error) {
          toast({
            title: t("admin.vision.load-error"),
            description: res.error,
            variant: "destructive",
          });
          // 即使出错也设置空值,让用户可以继续操作
          setModels("");
        }
      } catch (error) {
        console.error("Failed to load vision config:", error);
        toast({
          title: t("admin.vision.load-error"),
          description: String(error),
          variant: "destructive",
        });
        setModels("");
      } finally {
        setLoading(false);
      }
    }
    loadConfig();
  }, [t, toast]);

  // 保存配置
  const handleSave = async () => {
    setSaving(true);
    const modelList = models
      .split("\n")
      .map((m) => m.trim())
      .filter((m) => m.length > 0);

    const res = await updateVisionConfig(modelList);
    toastState(toast, t, res, true);
    setSaving(false);
  };

  // 刷新配置(使生效)
  const handleRefresh = async () => {
    setRefreshing(true);
    const res = await refreshVisionConfig();
    toastState(toast, t, res, true);
    setRefreshing(false);
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Loader2 className="w-8 h-8 animate-spin" />
      </div>
    );
  }

  return (
    <div className="vision-page p-4 max-w-4xl mx-auto">
      <Card>
        <CardHeader>
          <CardTitle>{t("admin.vision.title")}</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <Alert>
            <Info className="h-4 w-4" />
            <AlertDescription>
              {t("admin.vision.description")}
            </AlertDescription>
          </Alert>

          <div className="space-y-2">
            <label className="text-sm font-medium">
              {t("admin.vision.models-label")}
            </label>
            <Textarea
              value={models}
              onChange={(e) => setModels(e.target.value)}
              placeholder={t("admin.vision.models-placeholder")}
              rows={12}
              className="font-mono text-sm"
            />
            <p className="text-sm text-muted-foreground">
              {t("admin.vision.models-hint")}
            </p>
          </div>

          <div className="flex gap-3 pt-4">
            <Button
              onClick={handleSave}
              disabled={saving || refreshing}
              className="flex items-center gap-2"
            >
              {saving ? (
                <>
                  <Loader2 className="w-4 h-4 animate-spin" />
                  {t("admin.vision.saving")}
                </>
              ) : (
                <>
                  <Save className="w-4 h-4" />
                  {t("admin.vision.save")}
                </>
              )}
            </Button>

            <Button
              onClick={handleRefresh}
              disabled={saving || refreshing}
              variant="outline"
              className="flex items-center gap-2"
            >
              {refreshing ? (
                <>
                  <Loader2 className="w-4 h-4 animate-spin" />
                  {t("admin.vision.refreshing")}
                </>
              ) : (
                <>
                  <RefreshCw className="w-4 h-4" />
                  {t("admin.vision.refresh")}
                </>
              )}
            </Button>
          </div>

          <Alert className="mt-4">
            <Info className="h-4 w-4" />
            <AlertDescription>
              {t("admin.vision.refresh-hint")}
            </AlertDescription>
          </Alert>
        </CardContent>
      </Card>
    </div>
  );
}

export default Vision;
