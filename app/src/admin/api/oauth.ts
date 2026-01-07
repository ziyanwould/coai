import { CommonResponse } from "@/api/common.ts";
import axios from "axios";
import { getErrorMessage } from "@/utils/base.ts";

export type LinuxDoOAuth = {
  enabled: boolean;
  client_id: string;
  client_secret?: string;
  redirect_url: string;
  has_secret?: boolean;
};

export type OAuthConfig = {
  linux_do: LinuxDoOAuth;
};

export type OAuthResponse = CommonResponse & {
  data?: OAuthConfig;
};

export async function getOAuthConfig(): Promise<OAuthResponse> {
  try {
    const response = await axios.get("/admin/oauth/config");
    return response.data as OAuthResponse;
  } catch (e) {
    return { status: false, error: getErrorMessage(e) };
  }
}

export async function updateOAuthConfig(
  linuxDo: LinuxDoOAuth,
): Promise<CommonResponse> {
  try {
    const response = await axios.post("/admin/oauth/config", {
      linux_do: linuxDo,
    });
    return response.data as CommonResponse;
  } catch (e) {
    return { status: false, error: getErrorMessage(e) };
  }
}
