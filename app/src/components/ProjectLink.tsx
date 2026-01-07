/*
 * @Author: Liu Jiarong
 * @Date: 2025-08-05 19:47:37
 * @LastEditors: Liu Jiarong
 * @LastEditTime: 2025-08-05 21:04:04
 * @FilePath: /coai/app/src/components/ProjectLink.tsx
 * @Description: 
 * 
 * Copyright (c) 2025 by ${git_name_email}, All Rights Reserved. 
 */
import { Button } from "./ui/button.tsx";
import { useConversationActions, useMessages } from "@/store/chat.ts";
import { MessageSquarePlus } from "lucide-react";
// import Github from "@/components/ui/icons/Github.tsx";
// import Gitlab from "@/components/ui/icons/Gitlab.tsx";
// import { openWindow } from "@/utils/device.ts";

function ProjectLink() {
  const messages = useMessages();

  const { toggle } = useConversationActions();

  return messages.length > 0 ? (
    <Button
      variant="outline"
      size="icon"
      onClick={async () => await toggle(-1)}
    >
      <MessageSquarePlus className={`h-4 w-4`} />
    </Button>
  ) : (
    <>
    {/* <Button
      variant="outline"
      size="icon"
      onClick={() =>
        openWindow("https://github.com/ziyanwould/AIForAl")
      }
    >
      <Github className="h-[1.2rem] w-[1.2rem] rotate-0 scale-100" />
    </Button>
    <Button
      variant="outline"
      size="icon"
      title="Gitlab"
      onClick={() =>
        openWindow("https://git.liujiarong.top/explore")
      }
    >
      <Gitlab className="h-[1.2rem] w-[1.2rem] rotate-0 scale-100" />
    </Button> */}
    </>
  );
}

export default ProjectLink;
