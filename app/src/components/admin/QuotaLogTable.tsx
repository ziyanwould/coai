// app/src/components/admin/QuotaLogTable.tsx
import { useTranslation } from "react-i18next";
import { useToast } from "@/components/ui/use-toast.ts";
import { useState } from "react";
import {
  QuotaLogResponse,
  QuotaLogSortOption,
  getQuotaLogList,
} from "@/admin/api/chart.ts"; // 从你的 API 文件中导入
import { useEffectAsync } from "@/utils/hook.ts";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table.tsx";
import { Button } from "@/components/ui/button.tsx";
import { Loader2, Search } from "lucide-react";
import { Input } from "@/components/ui/input.tsx";
import { PaginationAction } from "@/components/ui/pagination.tsx";
import { isEnter } from "@/utils/base.ts";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select.tsx";

function QuotaLogTable() {
  const { t } = useTranslation();
  const { toast } = useToast();
  const [data, setData] = useState<QuotaLogResponse>({
    status: true,
    total: 0,
    data: [],
  });
  const [page, setPage] = useState<number>(0);
  const [search, setSearch] = useState<string>("");
  const [userIdFilter, setUserIdFilter] = useState<number | null>(null); // 新增用户ID筛选
  const [loading, setLoading] = useState<boolean>(false);
  const [sortKey, setSortKey] = useState<QuotaLogSortOption>("id"); // 默认按 ID 排序

  // 排序选项 (与后端 QuotaLogSortOption 类型对应)
  const sortOptions: { value: QuotaLogSortOption; label: string }[] = [
    { value: "id", label: "ID" },
    { value: "user_id", label: "User ID" },
    { value: "operation", label: "Operation" },
    { value: "quota_change", label: "Quota Change" },
    {value: "created_at", label: "Created At"},

  ];

  async function update() {
    setLoading(true);
    const resp = await getQuotaLogList(page, search, userIdFilter, sortKey);
    setLoading(false);
    if (resp.status) {
      setData(resp);
    } else {
      toast({
        title: t("admin.error"),
        description: resp.message,
      });
    }
  }

  useEffectAsync(update, [page, search, userIdFilter, sortKey]);

  return (
    <div className={`quota-log-table`}>
        <div className={`flex flex-row mb-6 items-center`}>
        <Input
          className={`search`}
          placeholder={t("admin.search-operation")} // 更改为搜索操作
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          onKeyDown={async (e) => {
            if (isEnter(e)) await update();
          }}
        />
        <Button size={`icon`} className={`flex-shrink-0 ml-2`} onClick={update}>
          <Search className={`h-4 w-4`} />
        </Button>
        {/* 用户 ID 筛选 */}
        <Input
              type="number"
              className={`w-[180px] ml-4`}
              placeholder={t("admin.filter-by-user-id")}
              value={userIdFilter === null ? "" : userIdFilter}
              onChange={(e) => {
                const value = e.target.value;
                setUserIdFilter(value === "" ? null : parseInt(value, 10));
              }}
              onKeyDown={async (e) => {
                if (isEnter(e)) await update();
              }}
        />
        {/*排序选项*/}
        <Select onValueChange={(value) => setSortKey(value as QuotaLogSortOption)} value={sortKey}>
          <SelectTrigger className="w-[180px] ml-4">
            <SelectValue placeholder={t("admin.sort-by")} />
          </SelectTrigger>
          <SelectContent>
            {sortOptions.map((option) => (
              <SelectItem key={option.value} value={option.value}>
                {t(`admin.${option.label.toLowerCase().replace(" ", "-")}`)} {/* 翻译标签 */}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>
      {(data.data && data.data.length > 0) || page > 0 ? (
        <>
          <Table>
            <TableHeader>
              <TableRow className={`select-none whitespace-nowrap`}>
                <TableHead>{t("admin.id")}</TableHead>
                <TableHead>{t("admin.user-id")}</TableHead>
                <TableHead>{t("admin.username")}</TableHead> 
                <TableHead>{t("admin.operation")}</TableHead>
                <TableHead>{t("admin.quota-change")}</TableHead>
                <TableHead>{t("admin.quota-before")}</TableHead>
                <TableHead>{t("admin.quota-after")}</TableHead>
                {/*已用配额的改变*/}
                 <TableHead>{t("admin.used-change")}</TableHead>
                <TableHead>{t("admin.used-before")}</TableHead>
                <TableHead>{t("admin.used-after")}</TableHead>
                <TableHead>{t("admin.created-at")}</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {data.data.map((log, idx) => (
                <TableRow key={idx}>
                  <TableCell>{log.id}</TableCell>
                  <TableCell>{log.user_id}</TableCell>
                  <TableCell>{log.username}</TableCell> 
                  <TableCell>{log.operation}</TableCell>
                  <TableCell>{log.quota_change}</TableCell>
                  <TableCell>{log.quota_before}</TableCell>
                  <TableCell>{log.quota_after}</TableCell>
                  <TableCell>{log.used_change}</TableCell>
                  <TableCell>{log.used_before}</TableCell>
                  <TableCell>{log.used_after}</TableCell>
                  <TableCell>{log.created_at}</TableCell>
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
          <p>{t("admin.no-quota-logs")}</p> {/* 添加一个空的提示 */}
        </div>
      )}
    </div>
  );
}

export default QuotaLogTable;