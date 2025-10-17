import { Play } from "lucide-react";
import { useState } from "react";

type MarkdownVideoProps = {
  src?: string;
  alt?: string;
};

export function MarkdownVideo({ src, alt }: MarkdownVideoProps) {
  const [isPlaying, setIsPlaying] = useState(false);

  if (!src) {
    return null;
  }

  // Check if it's a video URL (mp4, webm, mov, etc.)
  const videoExtensions = ['.mp4', '.webm', '.mov', '.avi', '.m4v', '.ogv'];
  const isVideo = videoExtensions.some(ext =>
    src.toLowerCase().includes(ext) ||
    src.toLowerCase().endsWith(ext)
  );

  if (!isVideo) {
    return null;
  }

  const handlePlay = () => {
    setIsPlaying(true);
  };

  return (
    <div className="markdown-video-wrapper relative group my-4">
      <div className="relative rounded-lg overflow-hidden bg-black">
        <video
          controls
          preload="metadata"
          className="w-full max-w-full h-auto"
          onPlay={handlePlay}
          poster="" // You can add a poster image if needed
        >
          <source src={src} type={getVideoMimeType(src)} />
          您的浏览器不支持视频播放。
        </video>

        {!isPlaying && (
          <div className="absolute inset-0 flex items-center justify-center bg-black bg-opacity-30 group-hover:bg-opacity-40 transition-all pointer-events-none">
            <div className="w-16 h-16 rounded-full bg-white bg-opacity-90 flex items-center justify-center shadow-lg">
              <Play className="w-8 h-8 text-black ml-1" />
            </div>
          </div>
        )}
      </div>

      {alt && (
        <div className="text-sm text-gray-500 dark:text-gray-400 mt-2 text-center">
          {alt}
        </div>
      )}
    </div>
  );
}

// Helper function to get video MIME type from URL
function getVideoMimeType(url: string): string {
  const lowerUrl = url.toLowerCase();

  if (lowerUrl.includes('.mp4') || lowerUrl.endsWith('.mp4')) {
    return 'video/mp4';
  }
  if (lowerUrl.includes('.webm') || lowerUrl.endsWith('.webm')) {
    return 'video/webm';
  }
  if (lowerUrl.includes('.mov') || lowerUrl.endsWith('.mov')) {
    return 'video/quicktime';
  }
  if (lowerUrl.includes('.avi') || lowerUrl.endsWith('.avi')) {
    return 'video/x-msvideo';
  }
  if (lowerUrl.includes('.m4v') || lowerUrl.endsWith('.m4v')) {
    return 'video/mp4';
  }
  if (lowerUrl.includes('.ogv') || lowerUrl.endsWith('.ogv')) {
    return 'video/ogg';
  }

  // Default to mp4
  return 'video/mp4';
}
