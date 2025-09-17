import { Download } from "lucide-react";
import { saveImageAsFile } from "@/utils/dom.ts";
import { Button } from "@/components/ui/button.tsx";

type MarkdownImageProps = {
  src?: string;
  alt?: string;
  enableDownload?: boolean;
};

export function MarkdownImage({ src, alt, enableDownload = false }: MarkdownImageProps) {
  if (!src) {
    return null;
  }

  const isBase64 = src.startsWith("data:image/");

  // Debug logging for image rendering
  if (isBase64) {
    console.log('[MarkdownImage] Rendering base64 image:', {
      srcLength: src.length,
      srcPrefix: src.substring(0, 50),
      alt,
      enableDownload
    });
  }

  // Generate filename from alt text or use default
  const getFilename = () => {
    if (alt && alt.trim() !== "" && alt !== "image") {
      // Clean alt text to make it a valid filename
      const cleanAlt = alt.replace(/[^a-zA-Z0-9\-_\u4e00-\u9fa5]/g, "_");
      return `${cleanAlt}.png`;
    }
    return `generated_image_${Date.now()}.png`;
  };

  const handleDownload = async () => {
    if (isBase64) {
      saveImageAsFile(getFilename(), src);
    } else {
      // For regular URL images, fetch and download
      try {
        const response = await fetch(src);
        const blob = await response.blob();
        const url = URL.createObjectURL(blob);
        const a = document.createElement("a");
        a.href = url;
        a.download = getFilename();
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        URL.revokeObjectURL(url);
      } catch (error) {
        console.error("Failed to download image:", error);
      }
    }
  };

  return (
    <div className="markdown-image-wrapper relative group">
      <img
        src={src}
        alt={alt || "Generated image"}
        className="markdown-image max-w-full h-auto rounded-lg"
      />
      {enableDownload && (
        <Button
          variant="secondary"
          size="sm"
          className="download-button"
          onClick={handleDownload}
          title="Download image"
        >
          <Download className="h-4 w-4 mr-1" />
          下载
        </Button>
      )}
    </div>
  );
}