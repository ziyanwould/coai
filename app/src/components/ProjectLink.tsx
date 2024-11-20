import { Button } from "./ui/button.tsx";
import { useConversationActions, useMessages } from "@/store/chat.ts";
import { MessageSquarePlus } from "lucide-react";
import Github from "@/components/ui/icons/Github.tsx";
import Gitlab from "@/components/ui/icons/Gitlab.tsx";
import { openWindow } from "@/utils/device.ts";

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
    <Button
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
    </Button>
    </>
  );
}

export default ProjectLink;
