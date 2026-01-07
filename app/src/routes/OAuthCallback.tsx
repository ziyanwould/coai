import { useToast } from "@/components/ui/use-toast.ts";
import { ToastAction } from "@/components/ui/toast.tsx";
import { tokenField } from "@/conf/bootstrap.ts";
import { useEffect } from "react";
import Loader from "@/components/Loader.tsx";
import "@/assets/pages/auth.less";
import { validateToken } from "@/store/auth.ts";
import { useDispatch } from "react-redux";
import router from "@/router.tsx";
import { useTranslation } from "react-i18next";
import { getQueryParam } from "@/utils/path.ts";
import { setMemory } from "@/utils/memory.ts";
import { goAuth } from "@/utils/app.ts";

function OAuthCallback() {
  const { toast } = useToast();
  const { t } = useTranslation();
  const dispatch = useDispatch();
  const token = getQueryParam("token").trim();
  const error = getQueryParam("error").trim();

  useEffect(() => {
    // 检查是否有错误
    if (error.length > 0) {
      const errorMessage = error.replace(/_/g, " ");
      toast({
        title: t("login-failed"),
        description: t("login-failed-prompt", { reason: errorMessage }),
        action: (
          <ToastAction altText={t("try-again")} onClick={goAuth}>
            {t("try-again")}
          </ToastAction>
        ),
      });

      setTimeout(goAuth, 3000);
      return;
    }

    // 检查是否有 token
    if (!token.length) {
      toast({
        title: t("invalid-token"),
        description: t("invalid-token-prompt"),
        action: (
          <ToastAction altText={t("try-again")} onClick={goAuth}>
            {t("try-again")}
          </ToastAction>
        ),
      });

      setTimeout(goAuth, 2500);
      return;
    }

    // 保存 token 并验证
    setMemory(tokenField, token);

    validateToken(dispatch, token, async () => {
      toast({
        title: t("login-success"),
        description: "欢迎回来!已通过 Linux.do 登录成功",
      });

      await router.navigate("/");
    });
  }, []);

  return (
    <div className={`auth`}>
      <Loader prompt={t("login")} />
    </div>
  );
}

export default OAuthCallback;
