// app/src/components/admin/UserTable.tsx
import { useTranslation } from "react-i18next";
import { useToast } from "@/components/ui/use-toast.ts";
import { useState } from "react";
import {
  CommonResponse,
  UserData,
  UserForm,
  UserResponse,
} from "@/admin/types.ts";
import {
  banUserOperation,
  getUserList,
  quotaOperation,
  releaseUsageOperation,
  setAdminOperation,
  subscriptionLevelOperation,
  subscriptionOperation,
  updateEmail,
  updatePassword,
} from "@/admin/api/chart.ts";
import { useEffectAsync } from "@/utils/hook.ts";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table.tsx";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu.tsx";
import { Button } from "@/components/ui/button.tsx";
import {
  CalendarCheck2,
  CalendarClock,
  CalendarOff,
  CloudCog,
  CloudFog,
  KeyRound,
  Loader2,
  Mail,
  MinusCircle,
  MoreHorizontal,
  PlusCircle,
  RotateCw,
  Search,
  Shield,
  ShieldMinus,
} from "lucide-react";
import { Input } from "@/components/ui/input.tsx";
import PopupDialog, { popupTypes } from "@/components/PopupDialog.tsx";
import { getNumber, isEnter, parseNumber } from "@/utils/base.ts";
import { useSelector } from "react-redux";
import { selectUsername } from "@/store/auth.ts";
import { PaginationAction } from "@/components/ui/pagination.tsx";
import Tips from "@/components/Tips.tsx";
import { Switch } from "@/components/ui/switch.tsx";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select.tsx";
import QuotaLogTable from "@/components/admin/QuotaLogTable"; // 导入 QuotaLogTable 组件

type OperationMenuProps = {
  user: UserData;
  onRefresh?: () => void;
};

export enum UserType {
  normal = "normal",
  basic_plan = "basic_plan",
  standard_plan = "standard_plan",
  pro_plan = "pro_plan",
}

export const userTypeArray = [
  UserType.normal,
  UserType.basic_plan,
  UserType.standard_plan,
  UserType.pro_plan,
];

function doToast(t: any, toast: any, resp: CommonResponse) {
  if (!resp.status)
    toast({
      title: t("admin.operate-failed"),
      description: t("admin.operate-failed-prompt", {
        reason: resp.message || resp.error,
      }),
    });
  else
    toast({
      title: t("admin.operate-success"),
      description: t("admin.operate-success-prompt"),
    });
}

function OperationMenu({ user, onRefresh }: OperationMenuProps) {
  const { t } = useTranslation();
  const { toast } = useToast();

  const username = useSelector(selectUsername);
  const [passwordOpen, setPasswordOpen] = useState<boolean>(false);
  const [emailOpen, setEmailOpen] = useState<boolean>(false);
  const [quotaOpen, setQuotaOpen] = useState<boolean>(false);
  const [quotaSetOpen, setQuotaSetOpen] = useState<boolean>(false);
  const [subscriptionOpen, setSubscriptionOpen] = useState<boolean>(false);
  const [subscriptionLevelOpen, setSubscriptionLevelOpen] =
    useState<boolean>(false);
  const [releaseOpen, setReleaseOpen] = useState<boolean>(false);
  const [banOpen, setBanOpen] = useState<boolean>(false);
  const [adminOpen, setAdminOpen] = useState<boolean>(false);

  return (
    <>
      <PopupDialog
        destructive={true}
        type={popupTypes.Text}
        title={t("admin.password-action")}
        name={t("auth.password")}
        description={t("admin.password-action-desc")}
        open={passwordOpen}
        setOpen={setPasswordOpen}
        defaultValue={""}
        onSubmit={async (password) => {
          const resp = await updatePassword(user.id, password);
          doToast(t, toast, resp);

          if (resp.status) {
            username === user.username && location.reload();
            onRefresh?.();
          }

          return resp.status;
        }}
      />
      <PopupDialog
        destructive={true}
        type={popupTypes.Text}
        title={t("admin.email-action")}
        name={t("admin.email")}
        description={t("admin.email-action-desc")}
        open={emailOpen}
        setOpen={setEmailOpen}
        defaultValue={user.email}
        onSubmit={async (email) => {
          const resp = await updateEmail(user.id, email);
          doToast(t, toast, resp);

          if (resp.status) onRefresh?.();
          return resp.status;
        }}
      />
      <PopupDialog
        type={popupTypes.Number}
        title={t("admin.quota-action")}
        name={t("admin.quota")}
        description={t("admin.quota-action-desc")}
        defaultValue={"0"}
        onValueChange={getNumber}
        open={quotaOpen}
        setOpen={setQuotaOpen}
        componentProps={{ acceptNegative: true }}
        onSubmit={async (value) => {
          const quota = parseNumber(value);
          const resp = await quotaOperation(user.id, quota);
          doToast(t, toast, resp);

          if (resp.status) onRefresh?.();
          return resp.status;
        }}
      />
      <PopupDialog
        type={popupTypes.Number}
        title={t("admin.quota-set-action")}
        name={t("admin.quota")}
        description={t("admin.quota-set-action-desc")}
        defaultValue={user.quota.toFixed(2)}
        onValueChange={getNumber}
        open={quotaSetOpen}
        setOpen={setQuotaSetOpen}
        componentProps={{ acceptNegative: true }}
        onSubmit={async (value) => {
          const quota = parseNumber(value);
          const resp = await quotaOperation(user.id, quota, true);
          doToast(t, toast, resp);

          if (resp.status) onRefresh?.();
          return resp.status;
        }}
      />
      <PopupDialog
        type={popupTypes.Clock}
        title={t("admin.subscription-action")}
        description={t("admin.subscription-action-desc", {
          username: user.username,
        })}
        defaultValue={user.expired_at || "1970-01-01 00:00:00"}
        open={subscriptionOpen}
        setOpen={setSubscriptionOpen}
        onSubmit={async (value) => {
          const resp = await subscriptionOperation(user.id, value);
          doToast(t, toast, resp);

          if (resp.status) onRefresh?.();
          return resp.status;
        }}
      />
      <PopupDialog
        type={popupTypes.List}
        title={t("admin.subscription-level")}
        name={t("admin.level")}
        description={t("admin.subscription-level-desc")}
        defaultValue={userTypeArray[user.level]}
        params={{
          dataList: userTypeArray,
          dataListTranslated: "admin.identity",
        }}
        open={subscriptionLevelOpen}
        setOpen={setSubscriptionLevelOpen}
        onSubmit={async (value) => {
          const level = userTypeArray.indexOf(value as UserType);
          console.log(level);
          const resp = await subscriptionLevelOperation(user.id, level);
          doToast(t, toast, resp);

          if (resp.status) onRefresh?.();
          return resp.status;
        }}
      />
      <PopupDialog
        type={popupTypes.Empty}
        title={t("admin.release-subscription-action")}
        name={t("admin.release-subscription")}
        description={t("admin.release-subscription-action-desc")}
        open={releaseOpen}
        setOpen={setReleaseOpen}
        onSubmit={async () => {
          const resp = await releaseUsageOperation(user.id);
          doToast(t, toast, resp);

          if (resp.status) onRefresh?.();
          return resp.status;
        }}
      />
      <PopupDialog
        disabled={username === user.username}
        destructive={true}
        type={popupTypes.Empty}
        title={user.is_banned ? t("admin.unban-action") : t("admin.ban-action")}
        description={
          user.is_banned
            ? t("admin.unban-action-desc")
            : t("admin.ban-action-desc")
        }
        open={banOpen}
        setOpen={setBanOpen}
        onSubmit={async () => {
          const resp = await banUserOperation(user.id, !user.is_banned);
          doToast(t, toast, resp);

          if (resp.status) onRefresh?.();
          return resp.status;
        }}
      />
      <PopupDialog
        disabled={username === user.username}
        destructive={true}
        type={popupTypes.Empty}
        title={
          user.is_admin
            ? t("admin.cancel-admin-action")
            : t("admin.set-admin-action")
        }
        description={
          user.is_admin
            ? t("admin.cancel-admin-action-desc")
            : t("admin.set-admin-action-desc")
        }
        open={adminOpen}
        setOpen={setAdminOpen}
        onSubmit={async () => {
          const resp = await setAdminOperation(user.id, !user.is_admin);
          doToast(t, toast, resp);

          if (resp.status) onRefresh?.();
          return resp.status;
        }}
      />

      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant={`outline`} size={`icon`}>
            <MoreHorizontal className={`h-4 w-4`} />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent className={`min-w-[8.75rem]`}>
          <DropdownMenuItem onClick={() => setPasswordOpen(true)}>
            <KeyRound className={`h-4 w-4 mr-2`} />
            {t("admin.password-action")}
          </DropdownMenuItem>
          <DropdownMenuItem onClick={() => setEmailOpen(true)}>
            <Mail className={`h-4 w-4 mr-2`} />
            {t("admin.email-action")}
          </DropdownMenuItem>
          {user.is_banned ? (
            <DropdownMenuItem onClick={() => setBanOpen(true)}>
              <PlusCircle className={`h-4 w-4 mr-2`} />
              {t("admin.unban-action")}
            </DropdownMenuItem>
          ) : (
            <DropdownMenuItem onClick={() => setBanOpen(true)}>
              <MinusCircle className={`h-4 w-4 mr-2`} />
              {t("admin.ban-action")}
            </DropdownMenuItem>
          )}
          {user.is_admin ? (
            <DropdownMenuItem onClick={() => setAdminOpen(true)}>
              <ShieldMinus className={`h-4 w-4 mr-2`} />
              {t("admin.cancel-admin-action")}
            </DropdownMenuItem>
          ) : (
            <DropdownMenuItem onClick={() => setAdminOpen(true)}>
              <Shield className={`h-4 w-4 mr-2`} />
              {t("admin.set-admin-action")}
            </DropdownMenuItem>
          )}
          <DropdownMenuItem onClick={() => setQuotaOpen(true)}>
            <CloudFog className={`h-4 w-4 mr-2`} />
            {t("admin.quota-action")}
          </DropdownMenuItem>
          <DropdownMenuItem onClick={() => setQuotaSetOpen(true)}>
            <CloudCog className={`h-4 w-4 mr-2`} />
            {t("admin.quota-set-action")}
          </DropdownMenuItem>
          <DropdownMenuItem onClick={() => setSubscriptionOpen(true)}>
            <CalendarClock className={`h-4 w-4 mr-2`} />
            {t("admin.subscription-action")}
          </DropdownMenuItem>
          <DropdownMenuItem onClick={() => setSubscriptionLevelOpen(true)}>
            <CalendarCheck2 className={`h-4 w-4 mr-2`} />
            {t("admin.subscription-level")}
          </DropdownMenuItem>
          <DropdownMenuItem onClick={() => setReleaseOpen(true)}>
            <CalendarOff className={`h-4 w-4 mr-2`} />
            {t("admin.release-subscription-action")}
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </>
  );
}

function UserTable() {
  const { t } = useTranslation();
  const { toast } = useToast();
  const [data, setData] = useState<UserForm>({
    total: 0,
    data: [],
  });
  const [page, setPage] = useState<number>(0);
  const [search, setSearch] = useState<string>("");
  const [isSubscribedFilter, setIsSubscribedFilter] = useState<boolean | null>(
    null,
  ); // 新增订阅状态筛选
  const [isBannedFilter, setIsBannedFilter] = useState<boolean | null>(null); // 新增封禁状态筛选
  const [loading, setLoading] = useState<boolean>(false);
  const [sortKey, setSortKey] = useState<SortOption>("id"); // 默认按 ID 排序

  type SortOption = "id" | "quota" | "used_quota" | "total_month" | "expired_at";

  const sortOptions: { value: SortOption; label: string }[] = [
    { value: "id", label: "ID" },
    { value: "quota", label: "Quota" },
    { value: "used_quota", label: "Used Quota" },
    { value: "total_month", label: "Total Month" },
    { value: "expired_at", label: "Expired At" },
  ];

  async function update() {
    setLoading(true);
    // 修改 getUserList 的调用，传递 emailSearch 和 isSubscribedFilter 和 isBannedFilter
    const resp = await getUserList(
      page,
      search,
      isSubscribedFilter,
      isBannedFilter,
      sortKey,
    );
    setLoading(false);
    if (resp.status) setData(resp as UserResponse);
    else
      toast({
        title: t("admin.error"),
        description: resp.message,
      });
  }
  useEffectAsync(update, [
    page,
    search,
    isSubscribedFilter,
    isBannedFilter,
    sortKey,
  ]); // 监听 sortKey 的变化 // 监听 search 和 isSubscribedFilter 和 isBannedFilter 的变化

  return (
    <div className={`user-table`}>
      <div className={`flex flex-row mb-6 items-center`}>
        {" "}
        {/* 使用 flex-row 和 items-center 居中对齐 */}
        <Input
          className={`search`}
          placeholder={t("admin.search-username")} // 修改 placeholder 提示可以搜索用户名和邮箱
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          onKeyDown={async (e) => {
            if (isEnter(e)) await update();
          }}
        />
        <Button
          size={`icon`}
          className={`flex-shrink-0 ml-2`}
          onClick={update}
        >
          <Search className={`h-4 w-4`} />
        </Button>
        <Select
          onValueChange={(value) => setSortKey(value as SortOption)}
          value={sortKey}
        >
          <SelectTrigger className="w-[180px] ml-4">
            <SelectValue placeholder={t("admin.sort-by")} />
          </SelectTrigger>
          <SelectContent>
            {sortOptions.map((option) => (
              <SelectItem key={option.value} value={option.value}>
                {t(`admin.${option.label.toLowerCase().replace(" ", "-")}`)}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
        <div className={`w-[200px] ml-4 flex items-center`}>
          {" "}
          {/*  添加容器包裹 Switch 和 label */}
          <Switch
            id="is-subscribed-filter"
            checked={isSubscribedFilter === true} //  当 isSubscribedFilter 为 true 时，Switch 选中
            onCheckedChange={(state) => {
              setIsSubscribedFilter(
                state === true ? true : state === false ? false : null,
              ); //  处理 null 情况，三元表达式更清晰
              setPage(0); //  切换筛选条件后，重置页码到第一页
            }}
          />
          <label
            htmlFor="is-subscribed-filter"
            className="ml-2 text-sm text-gray-500 dark:text-gray-400"
          >
            {t("admin.is-subscribed")}
          </label>{" "}
          {/* 添加 label，并使用 htmlFor 关联 Switch */}
        </div>
        <div className={`w-[180px] ml-4 flex items-center`}>
          {" "}
          {/*  添加容器包裹 Switch 和 label */}
          <Switch
            id="is-banned-filter"
            checked={isBannedFilter === true}
            onCheckedChange={(state) => {
              setIsBannedFilter(
                state === true ? true : state === false ? false : null,
              );
              setPage(0);
            }}
          />
          <label
            htmlFor="is-banned-filter"
            className="ml-2 text-sm text-gray-500 dark:text-gray-400"
          >
            {t("admin.is-banned")}
          </label>{" "}
          {/* 添加 label */}
        </div>
      </div>
      {(data.data && data.data.length > 0) || page > 0 ? (
        <>
          <Table>
            <TableHeader>
              <TableRow className={`select-none whitespace-nowrap`}>
                <TableHead>ID</TableHead>
                <TableHead>{t("admin.username")}</TableHead>
                <TableHead>{t("admin.email")}</TableHead>
                <TableHead>{t("admin.quota")}</TableHead>
                <TableHead>{t("admin.used-quota")}</TableHead>
                <TableHead>{t("admin.is-subscribed")}</TableHead>
                <TableHead>{t("admin.level")}</TableHead>
                <TableHead>{t("admin.total-month")}</TableHead>
                <TableHead>{t("admin.expired-at")}</TableHead>
                <TableHead>{t("admin.is-banned")}</TableHead>
                <TableHead>{t("admin.is-admin")}</TableHead>
                <TableHead>{t("admin.action")}</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {(data.data || []).map((user, idx) => (
                <TableRow key={idx}>
                  <TableCell>{user.id}</TableCell>
                  <TableCell className={`whitespace-nowrap`}>
                    {user.username}
                  </TableCell>
                  <TableCell className={`whitespace-nowrap`}>
                    {user.email || "-"}
                  </TableCell>
                  <TableCell>{user.quota}</TableCell>
                  <TableCell>{user.used_quota}</TableCell>
                  <TableCell>
                    {t(user.is_subscribed.toString())}
                    <Tips
                      className={`inline-block`}
                      content={t("admin.is-subscribed-tips")}
                    />
                  </TableCell>
                  <TableCell className={`whitespace-nowrap`}>
                    {t(`admin.identity.${userTypeArray[user.level]}`)}
                  </TableCell>
                  <TableCell>{user.total_month}</TableCell>
                  <TableCell className={`whitespace-nowrap`}>
                    {user.expired_at || "-"}
                  </TableCell>
                  <TableCell>{t(user.is_banned.toString())}</TableCell>
                  <TableCell>{t(user.is_admin.toString())}</TableCell>
                  <TableCell>
                    <OperationMenu user={user} onRefresh={update} />
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
          <PaginationAction
            current={page}
            total={data.total}
            onPageChange={setPage}
            offset
          />
        </>
      ) : loading ? (
        <div className={`flex flex-col mb-4 mt-12 items-center`}>
          <Loader2 className={`w-6 h-6 inline-block animate-spin`} />
        </div>
      ) : (
        <div className={`empty`}>
          <p>{t("admin.empty")}</p>
        </div>
      )}
      <div className={`user-action`}>
        <div className={`grow`} />
        <Button variant={`outline`} size={`icon`} onClick={update}>
          <RotateCw className={`h-4 w-4`} />
        </Button>
      </div>
      {/* 新增的 QuotaLogTable */}
      <h2 className="text-2xl font-bold mt-8 mb-4">{t("admin.quota-log")}</h2>
      <QuotaLogTable />
    </div>
  );
}

export default UserTable;