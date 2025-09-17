import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Paintbrush } from 'lucide-react';
import { InpaintingCanvas } from './InpaintingCanvas';

interface InpaintingTriggerProps {
  imageUrl: string;
  prompt: string;
}

export function InpaintingTrigger({ imageUrl, prompt }: InpaintingTriggerProps) {
  const [isCanvasOpen, setIsCanvasOpen] = useState(false);

  const handleConfirm = async (maskDataUrl: string) => {
    setIsCanvasOpen(false);

    // Create a new message with both original image and mask
    const newPrompt = `${prompt}

![原图](${imageUrl})

![遮罩图](${maskDataUrl})`;

    // Update the current input field with the new prompt
    const inputElement = document.getElementById('input') as HTMLTextAreaElement;
    if (inputElement) {
      inputElement.value = newPrompt;
      inputElement.focus();

      // Trigger input event to notify React
      const event = new Event('input', { bubbles: true });
      inputElement.dispatchEvent(event);
    }
  };

  const handleCancel = () => {
    setIsCanvasOpen(false);
  };

  return (
    <div className="inpainting-trigger-wrapper p-4 border border-blue-200 rounded-lg bg-blue-50">
      <div className="flex items-center gap-3">
        <Paintbrush className="h-6 w-6 text-blue-600" />
        <div className="flex-1">
          <h3 className="font-medium text-blue-900">图像修复模式</h3>
          <p className="text-sm text-blue-700">
            点击下方按钮，在图片上标记需要重新生成的区域
          </p>
        </div>
        <Button
          onClick={() => setIsCanvasOpen(true)}
          className="bg-blue-600 hover:bg-blue-700"
        >
          <Paintbrush className="h-4 w-4 mr-2" />
          开始涂鸦
        </Button>
      </div>

      <InpaintingCanvas
        imageUrl={imageUrl}
        isOpen={isCanvasOpen}
        onConfirm={handleConfirm}
        onCancel={handleCancel}
      />
    </div>
  );
}