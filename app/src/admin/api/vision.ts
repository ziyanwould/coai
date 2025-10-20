import { CommonResponse } from "@/api/common.ts";
import axios from "axios";
import { getErrorMessage } from "@/utils/base.ts";

export type VisionConfig = {
  models: string[];
};

export type VisionResponse = CommonResponse & {
  data?: VisionConfig;
};

export async function getVisionConfig(): Promise<VisionResponse> {
  try {
    const response = await axios.get("/admin/vision/config");
    return response.data as VisionResponse;
  } catch (e) {
    return { status: false, error: getErrorMessage(e) };
  }
}

export async function updateVisionConfig(
  models: string[],
): Promise<CommonResponse> {
  try {
    const response = await axios.post("/admin/vision/config", { models });
    return response.data as CommonResponse;
  } catch (e) {
    return { status: false, error: getErrorMessage(e) };
  }
}

export async function refreshVisionConfig(): Promise<CommonResponse> {
  try {
    const response = await axios.post("/admin/vision/refresh");
    return response.data as CommonResponse;
  } catch (e) {
    return { status: false, error: getErrorMessage(e) };
  }
}
