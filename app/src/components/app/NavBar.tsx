/*
 * @Author: Liu Jiarong
 * @Date: 2024-06-08 21:52:41
 * @LastEditors: Liu Jiarong
 * @LastEditTime: 2024-09-08 08:45:27
 * @FilePath: /chatnio/app/src/components/app/NavBar.tsx
 * @Description: 
 * 
 * Copyright (c) 2024 by ${git_name_email}, All Rights Reserved. 
 */
import "@/assets/pages/navbar.less";
import { useTranslation } from "react-i18next";
import { useDispatch, useSelector } from "react-redux";
import {
  selectAuthenticated,
  selectUsername,
  validateToken,
} from "@/store/auth.ts";
import { Button } from "@/components/ui/button.tsx";
import { Menu } from "lucide-react";
import { useEffect } from "react";
import { tokenField } from "@/conf/bootstrap.ts";
import { toggleMenu } from "@/store/menu.ts";
import ProjectLink from "@/components/ProjectLink.tsx";
import { ModeToggle } from "@/components/ThemeProvider.tsx";
import router from "@/router.tsx";
import MenuBar from "./MenuBar.tsx";
import { getMemory } from "@/utils/memory.ts";
import { goAuth } from "@/utils/app.ts";
import Avatar from "@/components/Avatar.tsx";
import { appLogo } from "@/conf/env.ts";
import Announcement from "@/components/app/Announcement.tsx";
import WebsiteIcon from "@/components/ui/icons/WebsiteIcon.tsx";
// import DoubleBubbleIcon from "@/components/ui/icons/DoubleBubbleIcon.tsx";
import ChatGPTIcon from "@/components/ui/icons/ChatGPTIcon.tsx";
import AIGenerationIcon from "@/components/ui/icons/AIGenerationIcon.tsx";
import DrawingIcon from "@/components/ui/icons/DrawingIcon.tsx";
import NavigationIcon from "@/components/ui/icons/NavigationIcon.tsx";
import { openWindow } from "@/utils/device.ts";

function NavMenu() {
  const username = useSelector(selectUsername);

  return (
    <div className={`avatar`}>
      <MenuBar>
        <Button variant={`ghost`} size={`icon`}>
          <Avatar username={username} />
        </Button>
      </MenuBar>
    </div>
  );
}

function NavBar() {
  const { t } = useTranslation();
  const dispatch = useDispatch();
  useEffect(() => {
    validateToken(dispatch, getMemory(tokenField));
  }, []);
  const auth = useSelector(selectAuthenticated);

  return (
    <nav className={`navbar`}>
      <div className={`items`}>
        <Button
          size={`icon`}
          variant={`ghost`}
          onClick={() => dispatch(toggleMenu())}
        >
          <Menu />
        </Button>
        <img
          className={`logo`}
          src={appLogo}
          alt=""
          onClick={() => router.navigate("/")}
        />
        <div className={`grow`} />
        <ProjectLink />
        <Button
          variant="outline"
          size="icon"
          title="AI-Chat国外大模型"
          onClick={() =>
            openWindow("https://robot.liujiarong.top")
          }
        >
          <ChatGPTIcon className="h-[1.2rem] w-[1.2rem] rotate-0 scale-100" />
        </Button>
        {/* <Button
          variant="outline"
          size="icon"
          title="国内加速版"
          onClick={() =>
            openWindow("https://robotai.liujiarong.top")
          }
        >
          <DoubleBubbleIcon className="h-[1.2rem] w-[1.2rem] rotate-0 scale-100" />
        </Button> */}
        <Button
          variant="outline"
          size="icon"
          title="首页"
          onClick={() =>
            openWindow("https://www.liujiarong.top")
          }
        >
          <WebsiteIcon className="h-[1.2rem] w-[1.2rem] rotate-0 scale-100" />
        </Button>
        <Button
          variant="outline"
          size="icon"
          title="阿里通义万相（AI生图）"
          onClick={() =>
            openWindow("https://www.liujiarong.top/wanx-ui")
          }
        >
          <AIGenerationIcon className="h-[1.2rem] w-[1.2rem] rotate-0 scale-100" />
        </Button>
        <Button
          variant="outline"
          size="icon"
          title="stable-ui（AI生图）"
          onClick={() =>
            openWindow("https://www.liujiarong.top/stable-ui/")
          }
        >
          <DrawingIcon className="h-[1.2rem] w-[1.2rem] rotate-0 scale-100" />
        </Button>
        <Button
          variant="outline"
          size="icon"
          title="网址导航"
          onClick={() =>
            openWindow("https://homarr.liujiarong.top/board")
          }
        >
          <NavigationIcon className="h-[1.2rem] w-[1.2rem] rotate-0 scale-100" />
        </Button>
        <Announcement />
        <ModeToggle />
        {auth ? (
          <NavMenu />
        ) : (
          <Button size={`sm`} onClick={goAuth}>
            {t("login")}
          </Button>
        )}
      </div>
    </nav>
  );
}

export default NavBar;
