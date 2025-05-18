import React, { FC } from "react";
import { Minus, Plus } from "lucide-react";

import { Button } from "components";

export type TZoomControlsProps = {
  zoomIn: () => void;
  zoomOut: () => void;
  scale: number;
};

export const ZoomControls: FC<TZoomControlsProps> = ({ zoomIn, zoomOut, scale }) => {
  return (
    <div className="absolute bottom-4 right-4 z-10 flex flex-col rounded-lg p-2 shadow-lg">
      <Button
        onClick={zoomIn}
        variant="ghost"
        aria-label="Zoom in"
      >
        <Plus className="h-5 w-5" />
      </Button>
      <div className="py-1 text-center font-mono text-sm text-gray-500">{Math.round(scale * 100)}%</div>
      <Button
        onClick={zoomOut}
        variant="ghost"
        aria-label="Zoom out"
      >
        <Minus className="h-5 w-5" />
      </Button>
    </div>
  );
};
