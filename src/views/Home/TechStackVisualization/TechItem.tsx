import React from "react";
import { motion } from "framer-motion";

import { TTechnology } from "./types";

type TTechItemProps = {
  tech: TTechnology;
  color: string;
};

export const TechItem: React.FC<TTechItemProps> = ({ tech, color }) => {
  return (
    <motion.div
      className="mb-3 flex flex-col items-center"
      whileHover={{ scale: 1.1 }}
      whileTap={{ scale: 0.95 }}
    >
      <motion.div
        className="mb-1 flex h-full w-full items-center justify-center"
        initial={{ rotate: 0 }}
        whileHover={{ rotate: 360, transition: { duration: 0.5 } }}
      >
        <TechIcon
          iconSrc={tech.icon}
          name={tech.name}
          color={color}
        />
      </motion.div>
      <motion.div
        className="text-center text-xs font-medium"
        initial={{ opacity: 0, y: 5 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.1 }}
      >
        {tech.name}
      </motion.div>
      <motion.div
        className="text-center text-xs text-gray-400"
        initial={{ opacity: 0, y: 5 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.2 }}
      >
        {tech.description}
      </motion.div>
    </motion.div>
  );
};

const TechIcon: React.FC<{ iconSrc?: string; name: string; color: string }> = ({ iconSrc, name, color }) => {
  return iconSrc ? (
    <img
      src={iconSrc}
      alt="MSW Logo"
      className="h-8 w-8"
    />
  ) : (
    <div
      className="flex h-8 w-8 items-center justify-center rounded-full text-xs font-bold text-white shadow-md"
      style={{ backgroundColor: color }}
    >
      {name.substring(0, 2)}
    </div>
  );
};
