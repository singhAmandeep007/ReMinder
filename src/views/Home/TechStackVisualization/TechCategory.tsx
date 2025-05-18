import React, { useMemo, useState } from "react";
import { motion, AnimatePresence } from "framer-motion";

import { TechItem } from "./TechItem";
import { TTechCategory } from "./types";

type TTechCategoryProps = {
  category: TTechCategory;
  position: { x: number; y: number };
  scale: number;
};

export const TechCategory: React.FC<TTechCategoryProps> = ({ category, position, scale }) => {
  const [isExpanded, setIsExpanded] = useState(false);
  const { name, color, icon, technologies } = category;

  const techPositions = useMemo(
    () =>
      technologies.map((_, index) => {
        const angle = (2 * Math.PI * index) / technologies.length;
        const radius = 100 * scale;
        const x = Math.cos(angle) * radius;
        const y = Math.sin(angle) * radius;
        return { x, y };
      }),
    [technologies, scale]
  );

  return (
    <motion.div
      className="absolute -translate-x-1/2 -translate-y-1/2 transform"
      style={{
        left: `${position.x}px`,
        top: `${position.y}px`,
      }}
      initial={false}
    >
      <div className="relative">
        {/* Glass effect backdrop for the main button */}
        <motion.div
          className="absolute inset-0 z-0 h-20 w-20 rounded-full"
          style={{
            background: `linear-gradient(135deg, ${color}cc, ${color}b3)`, // darker alpha
            backdropFilter: "blur(8px)",
            boxShadow: `0 8px 32px 0 rgba(31, 38, 135, 0.37)`,
          }}
          initial={false}
          animate={{
            scale: isExpanded ? 1.1 : 1,
          }}
        />

        <motion.div
          className="relative z-10 flex h-20 w-20 cursor-pointer flex-col items-center justify-center rounded-full p-2"
          whileHover={{ scale: 1.1 }}
          whileTap={{ scale: 0.95 }}
          onClick={() => setIsExpanded(!isExpanded)}
        >
          <motion.div
            className="mb-1 text-white"
            initial={false}
            animate={{
              rotate: isExpanded ? 360 : 0,
              scale: isExpanded ? 1.2 : 1,
            }}
            transition={{
              rotate: { duration: 0.5 },
              scale: { duration: 0.3 },
            }}
          >
            {icon}
          </motion.div>
          <div className="text-center text-xs font-bold text-white">{name}</div>
        </motion.div>

        {/* Glowing ring effect when expanded */}
        {isExpanded && (
          <motion.div
            className="absolute inset-0 -z-10 h-20 w-20 rounded-full"
            initial={{ opacity: 0, scale: 0.8 }}
            animate={{
              opacity: 1,
              scale: 1.3,
              boxShadow: `0 0 20px 5px ${color}80`,
            }}
            exit={{ opacity: 0, scale: 1 }}
          />
        )}

        <AnimatePresence>
          {isExpanded &&
            technologies.map((tech, index) => (
              <motion.div
                key={tech.id}
                className="absolute -translate-x-1/2 -translate-y-1/2 transform"
                style={{
                  left: techPositions[index].x + "px",
                  top: techPositions[index].y + "px",
                }}
                initial={{ scale: 0, opacity: 0 }}
                animate={{
                  scale: 1,
                  opacity: 1,
                  transition: {
                    delay: index * 0.05,
                    type: "spring",
                    stiffness: 260,
                    damping: 20,
                  },
                }}
                exit={{
                  scale: 0,
                  opacity: 0,
                  transition: {
                    duration: 0.2,
                    delay: (technologies.length - index) * 0.02,
                  },
                }}
              >
                <TechItem
                  tech={tech}
                  color={color}
                />
              </motion.div>
            ))}
        </AnimatePresence>
      </div>
    </motion.div>
  );
};
