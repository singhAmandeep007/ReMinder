import React, { useState, useEffect, useRef } from "react";
import { motion, AnimatePresence } from "framer-motion";

import { Typography } from "components";

import { ZoomControls } from "./ZoomControls";
import { TechCategory } from "./TechCategory";
import { ConnectionLines } from "./ConnectionLines";
import { techStack } from "./data";

export const TechStackVisualization: React.FC = () => {
  const [scale, setScale] = useState(1);
  const [showCategories, setShowCategories] = useState(true);
  const [center, setCenter] = useState<{ x: number; y: number } | null>(null);

  const containerRef = useRef<HTMLDivElement>(null);

  // Define category positions in a circular layout
  const getCategoryPositions = () => {
    if (!containerRef.current) return {};

    const categories = Object.keys(techStack);
    const rect = containerRef.current.getBoundingClientRect();
    const centerX = rect.width / 2;
    const centerY = rect.height / 2;
    const radius = Math.min(rect.width, rect.height) * 0.35;

    return categories.reduce(
      (acc, key, index) => {
        const angle = (2 * Math.PI * (index + 1.5)) / categories.length;
        const x = centerX + Math.cos(angle) * radius;
        const y = centerY + Math.sin(angle) * radius;
        acc[key] = { x, y };
        return acc;
      },
      {} as Record<string, { x: number; y: number }>
    );
  };

  const getCenter = () => {
    if (!containerRef.current) return null;

    const rect = containerRef.current.getBoundingClientRect();
    const centerX = rect.width / 2;
    const centerY = rect.height / 2;

    return { x: centerX, y: centerY };
  };

  const [categoryPositions, setCategoryPositions] = useState<Record<string, { x: number; y: number }>>({});

  // Update positions on container resize
  useEffect(() => {
    const updatePositions = () => {
      setCategoryPositions(getCategoryPositions());

      setCenter(getCenter());
    };

    const resizeObserver = new ResizeObserver(updatePositions);
    if (containerRef.current) {
      resizeObserver.observe(containerRef.current);
    }

    updatePositions();
    return () => resizeObserver.disconnect();
  }, []);

  const zoomIn = () => {
    setScale((prevScale) => Math.min(prevScale + 0.1, 2));
  };

  const zoomOut = () => {
    setScale((prevScale) => Math.max(prevScale - 0.1, 0.5));
  };

  const toggleCategories = () => {
    setShowCategories((prev) => !prev);
  };

  // Connection points for drawing lines
  const connectionPoints = Object.entries(techStack).map(([key, category]) => ({
    x: categoryPositions[key]?.x || 0,
    y: categoryPositions[key]?.y || 0,
    color: category.color,
  }));

  return (
    <div className="relative h-full w-full overflow-hidden rounded-xl text-foreground">
      <div className="absolute left-4 top-4 z-10 rounded-md px-3 py-1.5 text-sm backdrop-blur-sm">
        <Typography variant="h6">Tech Stack Architecture</Typography>
      </div>
      <div
        className="h-full w-full transition-transform duration-200"
        style={{ transform: `scale(${scale})` }}
        ref={containerRef}
      >
        <AnimatePresence>
          {showCategories && (
            <>
              {center && (
                <ConnectionLines
                  connections={connectionPoints}
                  scale={2}
                  center={center}
                />
              )}

              {Object.entries(techStack).map(([key, category]) => (
                <motion.div
                  key={key}
                  initial={{ opacity: 0 }}
                  animate={{ scale: 1, opacity: 1 }}
                  exit={{ opacity: 0 }}
                  transition={{ duration: 0.5 }}
                >
                  <TechCategory
                    category={category}
                    position={categoryPositions[key] || { x: 0, y: 0 }}
                    scale={1}
                  />
                </motion.div>
              ))}
            </>
          )}
        </AnimatePresence>

        {/* Central icon */}
        <div className="absolute left-1/2 top-1/2 z-20 -translate-x-1/2 -translate-y-1/2 transform cursor-pointer">
          <motion.div
            onClick={toggleCategories}
            whileHover={{ scale: 1.1 }}
            whileTap={{ scale: 0.95 }}
          >
            <div className="flex h-32 w-32 flex-col items-center justify-center rounded-full border-4 border-gray-200 bg-white shadow-xl">
              <div className="mb-1 text-lg font-bold text-gray-800">ReMinder</div>
              <div className="px-2 text-center text-xs text-gray-500">Monorepo Architecture</div>
            </div>
          </motion.div>
        </div>
      </div>

      <ZoomControls
        zoomIn={zoomIn}
        zoomOut={zoomOut}
        scale={scale}
      />
    </div>
  );
};
