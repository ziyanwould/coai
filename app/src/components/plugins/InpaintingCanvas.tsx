import React, { useRef, useEffect, useState, useCallback } from 'react';
import { Button } from '@/components/ui/button';
import { Slider } from '@/components/ui/slider';
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Eraser, Paintbrush, RotateCcw, Check, X } from 'lucide-react';

interface InpaintingCanvasProps {
  imageUrl: string;
  isOpen: boolean;
  onConfirm: (maskDataUrl: string) => void;
  onCancel: () => void;
}

export function InpaintingCanvas({ imageUrl, isOpen, onConfirm, onCancel }: InpaintingCanvasProps) {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const [isDrawing, setIsDrawing] = useState(false);
  const [brushSize, setBrushSize] = useState([20]);
  const [isEraser, setIsEraser] = useState(false);
  const [imageLoaded, setImageLoaded] = useState(false);

  // Load and draw the image on canvas
  useEffect(() => {
    if (!isOpen || !imageUrl) return;

    const canvas = canvasRef.current;
    if (!canvas) return;

    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    const img = new Image();
    img.onload = () => {
      // Set canvas size to image size
      canvas.width = img.width;
      canvas.height = img.height;

      // Draw the image
      ctx.drawImage(img, 0, 0);

      // Set up mask drawing context
      ctx.globalCompositeOperation = 'source-over';
      ctx.lineCap = 'round';
      ctx.lineJoin = 'round';

      setImageLoaded(true);
    };
    img.src = imageUrl;
  }, [imageUrl, isOpen]);

  const startDrawing = useCallback((e: React.MouseEvent<HTMLCanvasElement>) => {
    const canvas = canvasRef.current;
    if (!canvas || !imageLoaded) return;

    const rect = canvas.getBoundingClientRect();
    const scaleX = canvas.width / rect.width;
    const scaleY = canvas.height / rect.height;

    const x = (e.clientX - rect.left) * scaleX;
    const y = (e.clientY - rect.top) * scaleY;

    setIsDrawing(true);

    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    ctx.globalCompositeOperation = isEraser ? 'destination-out' : 'source-over';
    ctx.strokeStyle = isEraser ? 'rgba(0,0,0,1)' : 'rgba(255,255,255,0.7)';
    ctx.lineWidth = brushSize[0];

    ctx.beginPath();
    ctx.moveTo(x, y);
  }, [brushSize, isEraser, imageLoaded]);

  const draw = useCallback((e: React.MouseEvent<HTMLCanvasElement>) => {
    if (!isDrawing) return;

    const canvas = canvasRef.current;
    if (!canvas) return;

    const rect = canvas.getBoundingClientRect();
    const scaleX = canvas.width / rect.width;
    const scaleY = canvas.height / rect.height;

    const x = (e.clientX - rect.left) * scaleX;
    const y = (e.clientY - rect.top) * scaleY;

    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    ctx.lineTo(x, y);
    ctx.stroke();
  }, [isDrawing]);

  const stopDrawing = useCallback(() => {
    setIsDrawing(false);
  }, []);

  const clearCanvas = () => {
    const canvas = canvasRef.current;
    if (!canvas) return;

    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    // Reload the original image
    const img = new Image();
    img.onload = () => {
      ctx.clearRect(0, 0, canvas.width, canvas.height);
      ctx.drawImage(img, 0, 0);
    };
    img.src = imageUrl;
  };

  const generateMask = () => {
    const canvas = canvasRef.current;
    if (!canvas) return;

    // Create a new canvas for the mask
    const maskCanvas = document.createElement('canvas');
    maskCanvas.width = canvas.width;
    maskCanvas.height = canvas.height;
    const maskCtx = maskCanvas.getContext('2d');
    if (!maskCtx) return;

    // Fill with black background
    maskCtx.fillStyle = 'black';
    maskCtx.fillRect(0, 0, maskCanvas.width, maskCanvas.height);

    // Get the current canvas image data
    const imageData = canvas.getContext('2d')?.getImageData(0, 0, canvas.width, canvas.height);
    if (!imageData) return;

    // Create mask data
    const maskData = maskCtx.createImageData(canvas.width, canvas.height);

    // Convert semi-transparent white areas to white, everything else to black
    for (let i = 0; i < imageData.data.length; i += 4) {
      const r = imageData.data[i];
      const g = imageData.data[i + 1];
      const b = imageData.data[i + 2];
      const a = imageData.data[i + 3];

      // If pixel is whitish with some transparency (user drawn area)
      if (r > 200 && g > 200 && b > 200 && a > 100) {
        maskData.data[i] = 255;     // R
        maskData.data[i + 1] = 255; // G
        maskData.data[i + 2] = 255; // B
        maskData.data[i + 3] = 255; // A
      } else {
        maskData.data[i] = 0;       // R
        maskData.data[i + 1] = 0;   // G
        maskData.data[i + 2] = 0;   // B
        maskData.data[i + 3] = 255; // A
      }
    }

    maskCtx.putImageData(maskData, 0, 0);

    // Convert to data URL
    const maskDataUrl = maskCanvas.toDataURL('image/png');
    onConfirm(maskDataUrl);
  };

  return (
    <Dialog open={isOpen} onOpenChange={onCancel}>
      <DialogContent className="max-w-4xl max-h-[90vh] overflow-auto">
        <DialogHeader>
          <DialogTitle>标记要修复的区域</DialogTitle>
        </DialogHeader>

        <div className="space-y-4">
          {/* Toolbar */}
          <div className="flex items-center gap-4 p-4 bg-gray-50 rounded-lg">
            <Button
              variant={!isEraser ? "default" : "outline"}
              size="sm"
              onClick={() => setIsEraser(false)}
            >
              <Paintbrush className="h-4 w-4 mr-2" />
              画笔
            </Button>

            <Button
              variant={isEraser ? "default" : "outline"}
              size="sm"
              onClick={() => setIsEraser(true)}
            >
              <Eraser className="h-4 w-4 mr-2" />
              橡皮擦
            </Button>

            <div className="flex items-center gap-2 ml-4">
              <span className="text-sm">笔刷大小:</span>
              <Slider
                value={brushSize}
                onValueChange={setBrushSize}
                max={50}
                min={5}
                step={1}
                className="w-32"
              />
              <span className="text-sm w-8">{brushSize[0]}</span>
            </div>

            <Button
              variant="outline"
              size="sm"
              onClick={clearCanvas}
              className="ml-auto"
            >
              <RotateCcw className="h-4 w-4 mr-2" />
              重置
            </Button>
          </div>

          {/* Canvas */}
          <div className="border rounded-lg overflow-hidden bg-gray-100">
            <canvas
              ref={canvasRef}
              className="max-w-full max-h-[60vh] object-contain cursor-crosshair"
              onMouseDown={startDrawing}
              onMouseMove={draw}
              onMouseUp={stopDrawing}
              onMouseLeave={stopDrawing}
              style={{
                display: imageLoaded ? 'block' : 'none',
                width: '100%',
                height: 'auto'
              }}
            />
            {!imageLoaded && (
              <div className="flex items-center justify-center h-64">
                <div className="text-gray-500">加载图片中...</div>
              </div>
            )}
          </div>

          {/* Instructions */}
          <div className="text-sm text-gray-600 bg-blue-50 p-3 rounded">
            <p><strong>使用说明：</strong></p>
            <p>• 使用画笔在图片上标记要修复的区域（白色区域将被重新生成）</p>
            <p>• 使用橡皮擦可以擦除标记</p>
            <p>• 调整笔刷大小来精确控制标记区域</p>
          </div>

          {/* Action buttons */}
          <div className="flex gap-2 justify-end">
            <Button variant="outline" onClick={onCancel}>
              <X className="h-4 w-4 mr-2" />
              取消
            </Button>
            <Button onClick={generateMask} disabled={!imageLoaded}>
              <Check className="h-4 w-4 mr-2" />
              确认并生成
            </Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}