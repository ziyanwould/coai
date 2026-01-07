import React from "react";
import { Codepen, Codesandbox, Github, Twitter, Youtube } from "lucide-react";
import { VirtualMessage } from "./VirtualMessage";
import { MarkdownVideo } from "@/components/plugins/video.tsx";

function getSocialIcon(url: string) {
  try {
    const { hostname } = new URL(url);

    if (hostname.includes("github.com"))
      return <Github className="h-4 w-4 inline-block mr-0.5" />;
    if (hostname.includes("twitter.com"))
      return <Twitter className="h-4 w-4 inline-block mr-0.5" />;
    if (hostname.includes("youtube.com"))
      return <Youtube className="h-4 w-4 inline-block mr-0.5" />;
    if (hostname.includes("codepen.io"))
      return <Codepen className="h-4 w-4 inline-block mr-0.5" />;
    if (hostname.includes("codesandbox.io"))
      return <Codesandbox className="h-4 w-4 inline-block mr-0.5" />;
  } catch (e) {
    return;
  }
}

type LinkProps = {
  href?: string;
  children: React.ReactNode;
};

export default function ({ href, children }: LinkProps) {
  const url: string = href?.toString() || "";

  if (url.startsWith("https://coai.virtual/reference::")) {
    const referenceUrl = url.slice("https://coai.virtual/reference::".length);
    return <VirtualMessage message={`reference::${referenceUrl}`}>{children}</VirtualMessage>;
  }

  if (url.startsWith("https://coai.virtual")) {
    const message = url.slice(20);

    return <VirtualMessage message={message}>{children}</VirtualMessage>;
  }

  // Check if this is a video URL
  const videoExtensions = ['.mp4', '.webm', '.mov', '.avi', '.m4v', '.ogv'];
  const isVideo = videoExtensions.some(ext =>
    url.toLowerCase().includes(ext) || url.toLowerCase().endsWith(ext)
  );

  if (isVideo) {
    // Return both the link and the video player below it
    return (
      <div className="video-link-wrapper">
        <a href={url} target={`_blank`} rel={`noopener noreferrer`} className="block mb-2">
          {getSocialIcon(url)}
          {children}
        </a>
        <MarkdownVideo src={url} alt={children?.toString()} />
      </div>
    );
  }

  return (
    <a href={url} target={`_blank`} rel={`noopener noreferrer`}>
      {getSocialIcon(url)}
      {children}
    </a>
  );
}
