/*
 * @Author: Liu Jiarong
 * @Date: 2024-06-10 23:58:33
 * @LastEditors: Liu Jiarong
 * @LastEditTime: 2024-06-10 23:58:51
 * @FilePath: /chatnio/app/src/components/ui/icons/NavigationIcon.tsx
 * @Description: 
 * 
 * Copyright (c) 2024 by ${git_name_email}, All Rights Reserved. 
 */
import React from "react";

function NavigationIcon(props: React.SVGProps<SVGSVGElement>) {
  return (
    <svg viewBox="0 0 1024 1024" {...props}>
      <path
        fill="#333333"
        d="M512 463a96 96 0 1 1 96-96 96.2 96.2 0 0 1-96 96z m0-128a32 32 0 1 0 32 32 32 32 0 0 0-32-32z"
      ></path>
      <path
        fill="#333333"
        d="M512 742a32.1 32.1 0 0 1-24.9-12l-170-211.1a250.7 250.7 0 1 1 390 0L536.9 730a32.1 32.1 0 0 1-24.9 12z m0-567.2a186.8 186.8 0 0 0-145.3 304.1l145.3 180 145.3-180A186.9 186.9 0 0 0 512 174.8z"
      ></path>
      <path
        fill="#333333"
        d="M810 913H214.5a86.2 86.2 0 0 1-79-121.2l60.5-135.4a86.7 86.7 0 0 1 79-51.1h46.2a32 32 0 1 1 0 64H275a22.5 22.5 0 0 0-20.5 13.2l-60.6 135.4a21.7 21.7 0 0 0 1.7 20.9 22.1 22.1 0 0 0 18.9 10H810a22.1 22.1 0 0 0 18.9-10 21.7 21.7 0 0 0 1.1-20.9l-60-135.3a22.4 22.4 0 0 0-20-13.3h-48.2a32 32 0 1 1 0-64h47.2a86.7 86.7 0 0 1 79 51.1l60 135.4A86.2 86.2 0 0 1 810 913z"
      ></path>
    </svg>
  );
}

export default NavigationIcon;