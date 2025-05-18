import React from "react";
import { motion } from "framer-motion";

type TConnectionPoint = {
  x: number;
  y: number;
  color: string;
};

type TConnectionLinesProps = {
  connections: TConnectionPoint[];
  scale: number;
  center: { x: number; y: number };
};

export const ConnectionLines: React.FC<TConnectionLinesProps> = ({ connections, scale, center }) => {
  if (connections.length < 2) return null;

  return (
    <svg className="pointer-events-none absolute left-0 top-0 h-full w-full">
      <g>
        {connections.map((start, i) => (
          <motion.line
            key={`center-${i}`}
            x2={start.x}
            y2={start.y}
            x1={center.x}
            y1={center.y}
            stroke={start.color}
            strokeWidth={2 * scale}
            strokeOpacity={0.5}
            strokeDasharray="5,5"
            initial={{ pathLength: 0, opacity: 0 }}
            animate={{
              pathLength: 1,
              opacity: 0.5,
              transition: {
                duration: 1,
                delay: i * 0.2,
                ease: "backOut",
              },
            }}
            exit={{
              opacity: 0,
              transition: {
                delay: 0,
              },
            }}
          />
        ))}
      </g>
    </svg>
  );
};
