import {
  BillingChartResponse,
  CommonResponse,
  ErrorChartResponse,
  InfoResponse,
  InvitationGenerateResponse,
  InvitationResponse,
  ModelChartResponse,
  RedeemResponse,
  RequestChartResponse,
  UserResponse,
  UserTypeChartResponse,
} from "@/admin/types.ts";
import axios from "axios";
import { getErrorMessage } from "@/utils/base.ts";

export async function getAdminInfo(): Promise<InfoResponse> {
  try {
    const response = await axios.get("/admin/analytics/info");
    return response.data as InfoResponse;
  } catch (e) {
    console.warn(e);
    return { subscription_count: 0, billing_today: 0, billing_month: 0 };
  }
}

export async function getModelChart(): Promise<ModelChartResponse> {
  try {
    const response = await axios.get("/admin/analytics/model");
    const data = response.data as ModelChartResponse;

    return {
      date: data.date,
      value: data.value || [],
    };
  } catch (e) {
    console.warn(e);
    return { date: [], value: [] };
  }
}

export async function getRequestChart(): Promise<RequestChartResponse> {
  try {
    const response = await axios.get("/admin/analytics/request");
    return response.data as RequestChartResponse;
  } catch (e) {
    console.warn(e);
    return { date: [], value: [] };
  }
}

export async function getBillingChart(): Promise<BillingChartResponse> {
  try {
    const response = await axios.get("/admin/analytics/billing");
    return response.data as BillingChartResponse;
  } catch (e) {
    console.warn(e);
    return { date: [], value: [] };
  }
}

export async function getErrorChart(): Promise<ErrorChartResponse> {
  try {
    const response = await axios.get("/admin/analytics/error");
    return response.data as ErrorChartResponse;
  } catch (e) {
    console.warn(e);
    return { date: [], value: [] };
  }
}

export async function getUserTypeChart(): Promise<UserTypeChartResponse> {
  try {
    const response = await axios.get("/admin/analytics/user");
    return response.data as UserTypeChartResponse;
  } catch (e) {
    console.warn(e);
    return {
      total: 0,
      normal: 0,
      api_paid: 0,
      basic_plan: 0,
      standard_plan: 0,
      pro_plan: 0,
    };
  }
}

export async function getInvitationList(
  page: number,
): Promise<InvitationResponse> {
  try {
    const response = await axios.get(`/admin/invitation/list?page=${page}`);
    return response.data as InvitationResponse;
  } catch (e) {
    return {
      status: false,
      message: getErrorMessage(e),
      data: [],
      total: 0,
    };
  }
}

export async function deleteInvitation(code: string): Promise<CommonResponse> {
  try {
    const response = await axios.post("/admin/invitation/delete", { code });
    return response.data as CommonResponse;
  } catch (e) {
    return { status: false, message: getErrorMessage(e) };
  }
}

export async function generateInvitation(
  type: string,
  quota: number,
  number: number,
): Promise<InvitationGenerateResponse> {
  try {
    const response = await axios.post("/admin/invitation/generate", {
      type,
      quota,
      number,
    });
    return response.data as InvitationGenerateResponse;
  } catch (e) {
    return { status: false, data: [], message: getErrorMessage(e) };
  }
}

export async function getRedeemList(page: number): Promise<RedeemResponse> {
  try {
    const response = await axios.get(`/admin/redeem/list?page=${page}`);
    return response.data as RedeemResponse;
  } catch (e) {
    console.warn(e);
    return { status: false, message: getErrorMessage(e), data: [], total: 0 };
  }
}

export async function deleteRedeem(code: string): Promise<CommonResponse> {
  try {
    const response = await axios.post("/admin/redeem/delete", { code });
    return response.data as CommonResponse;
  } catch (e) {
    return { status: false, message: getErrorMessage(e) };
  }
}

export async function generateRedeem(
  quota: number,
  number: number,
): Promise<InvitationGenerateResponse> {
  try {
    const response = await axios.post("/admin/redeem/generate", {
      quota,
      number,
    });
    return response.data as InvitationGenerateResponse;
  } catch (e) {
    return { status: false, data: [], message: getErrorMessage(e) };
  }
}
type SortOption = "id" | "quota" | "used_quota" | "total_month" | "expired_at";
export async function getUserList(
  page: number,
  search: string,
  isSubscribedFilter: boolean | null,
  isBannedFilter: boolean | null,
  sortKey: SortOption = "id" // 添加 sortKey 参数，默认为 "id"
): Promise<UserResponse> {
  try {
    let queryString = `/admin/user/list?page=${page}&search=${search}&sort=${sortKey}`; // 添加 sort 参数

    // 根据 isSubscribedFilter 的值，动态添加 is_subscribed query 参数
    if (isSubscribedFilter === true) {
      queryString += `&is_subscribed=true`;
    } else if (isSubscribedFilter === false) {
      queryString += `&is_subscribed=false`;
    }
    // 如果 isSubscribedFilter 为 null 或 undefined，则不添加 is_subscribed 参数，表示不进行订阅状态筛选

    // 根据 isBannedFilter 的值，动态添加 is_banned query 参数
    if (isBannedFilter === true) {
      queryString += `&is_banned=true`;
    } else if (isBannedFilter === false) {
      queryString += `&is_banned=false`;
    }
    // 如果 isBannedFilter 为 null 或 undefined，则不添加 is_banned 参数，表示不进行封禁状态筛选

    const response = await axios.get(queryString);
    return response.data as UserResponse;
  } catch (e) {
    return {
      status: false,
      message: getErrorMessage(e),
      data: [],
      total: 0,
    };
  }
}

export async function updatePassword(
  id: number,
  password: string,
): Promise<CommonResponse> {
  try {
    const response = await axios.post("/admin/user/password", {
      id,
      password,
    });
    return response.data as CommonResponse;
  } catch (e) {
    return { status: false, message: getErrorMessage(e) };
  }
}

export async function updateEmail(
  id: number,
  email: string,
): Promise<CommonResponse> {
  try {
    const response = await axios.post("/admin/user/email", {
      id,
      email,
    });
    return response.data as CommonResponse;
  } catch (e) {
    return { status: false, message: getErrorMessage(e) };
  }
}

export async function quotaOperation(
  id: number,
  quota: number,
  override?: boolean,
): Promise<CommonResponse> {
  try {
    const response = await axios.post("/admin/user/quota", {
      id,
      quota,
      override: override ?? false,
    });
    return response.data as CommonResponse;
  } catch (e) {
    return { status: false, message: getErrorMessage(e) };
  }
}

export async function subscriptionOperation(
  id: number,
  expired: string,
): Promise<CommonResponse> {
  try {
    const response = await axios.post("/admin/user/subscription", {
      id,
      expired,
    });
    return response.data as CommonResponse;
  } catch (e) {
    return { status: false, message: getErrorMessage(e) };
  }
}

export async function banUserOperation(
  id: number,
  ban: boolean,
): Promise<CommonResponse> {
  try {
    const response = await axios.post("/admin/user/ban", {
      id,
      ban,
    });
    return response.data as CommonResponse;
  } catch (e) {
    return { status: false, message: getErrorMessage(e) };
  }
}

export async function setAdminOperation(
  id: number,
  admin: boolean,
): Promise<CommonResponse> {
  try {
    const response = await axios.post("/admin/user/admin", {
      id,
      admin,
    });
    return response.data as CommonResponse;
  } catch (e) {
    return { status: false, message: getErrorMessage(e) };
  }
}

export async function subscriptionLevelOperation(
  id: number,
  level: number,
): Promise<CommonResponse> {
  try {
    const response = await axios.post("/admin/user/level", { id, level });
    return response.data as CommonResponse;
  } catch (e) {
    return { status: false, message: getErrorMessage(e) };
  }
}

export async function releaseUsageOperation(
  id: number,
): Promise<CommonResponse> {
  try {
    const response = await axios.post("/admin/user/release", { id });
    return response.data as CommonResponse;
  } catch (e) {
    return { status: false, message: getErrorMessage(e) };
  }
}
