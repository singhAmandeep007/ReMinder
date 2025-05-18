import { useState } from "react";
import { motion } from "framer-motion";
import { Palette, Wrench, Package, Zap, Layers } from "lucide-react";

export function DependencyVisualization() {
  const [activeCategory, setActiveCategory] = useState<string | null>(null);

  // Define dependency categories
  const categories = [
    {
      id: "development",
      name: "Development Tools",
      icon: <Wrench className="h-6 w-6" />,
      color: "#f97316",
      dependencies: [
        { name: "ESLint", description: "Code linting" },
        { name: "Prettier", description: "Code formatting" },
        { name: "Husky", description: "Git hooks" },
        { name: "lint-staged", description: "Staged file linting" },
      ],
    },
    {
      id: "styling",
      name: "Styling",
      icon: <Palette className="h-6 w-6" />,
      color: "#8b5cf6",
      dependencies: [
        { name: "Tailwind CSS", description: "Utility-first CSS" },
        { name: "Shadcn UI", description: "UI component library" },
        { name: "PostCSS", description: "CSS processing" },
        { name: "Autoprefixer", description: "CSS vendor prefixing" },
      ],
    },
    {
      id: "core",
      name: "Core",
      icon: <Layers className="h-6 w-6" />,
      color: "#22c55e",
      dependencies: [
        { name: "TypeScript", description: "Type-safe JavaScript" },
        { name: "React", description: "UI library" },
        { name: "React Router", description: "Routing" },
      ],
    },
    {
      id: "state",
      name: "Data Fetching & State",
      icon: <Zap className="h-6 w-6" />,
      color: "#3b82f6",
      dependencies: [
        { name: "Redux Toolkit", description: "State management" },
        { name: "RTK Query", description: "Data fetching" },
      ],
    },
    {
      id: "bundler",
      name: "Bundler",
      icon: <Package className="h-6 w-6" />,
      color: "#ef4444",
      dependencies: [
        { name: "Webpack", description: "Module bundler" },
        { name: "Babel", description: "JavaScript compiler" },
      ],
    },
  ];

  // Animation variants
  const containerVariants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: {
        staggerChildren: 0.1,
        delayChildren: 0.3,
      },
    },
  };

  const itemVariants = {
    hidden: { y: 20, opacity: 0 },
    visible: {
      y: 0,
      opacity: 1,
      transition: { type: "spring", stiffness: 300, damping: 24 },
    },
  };

  const dependencyVariants = {
    hidden: { scale: 0.8, opacity: 0 },
    visible: {
      scale: 1,
      opacity: 1,
      transition: {
        type: "spring",
        stiffness: 300,
        damping: 24,
        staggerChildren: 0.07,
      },
    },
  };

  const dependencyItemVariants = {
    hidden: { x: -20, opacity: 0 },
    visible: {
      x: 0,
      opacity: 1,
      transition: { type: "spring", stiffness: 300, damping: 24 },
    },
  };

  return (
    <div className="flex flex-col space-y-8">
      <motion.div
        className="grid grid-cols-1 gap-4 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-5"
        variants={containerVariants}
        initial="hidden"
        animate="visible"
      >
        {categories.map((category) => (
          <motion.div
            key={category.id}
            variants={itemVariants}
            className={`cursor-pointer rounded-lg border-2 p-4 transition-all duration-300 ${
              activeCategory === category.id
                ? `border-[${category.color}] bg-[${category.color}]/10`
                : `hover:border-[${category.color}]/50 border-border`
            }`}
            onClick={() => setActiveCategory(activeCategory === category.id ? null : category.id)}
            style={{
              borderColor: activeCategory === category.id ? category.color : undefined,
              backgroundColor: activeCategory === category.id ? `${category.color}10` : undefined,
            }}
          >
            <div className="flex h-full w-full justify-start gap-2">
              <div
                className="flex h-10 w-10 items-center justify-center rounded-full p-2"
                style={{ backgroundColor: `${category.color}20`, color: category.color }}
              >
                {category.icon}
              </div>
              <div className="self-start">
                <h3 className="font-semibold">{category.name}</h3>
                <p className="text-xs text-muted-foreground">{category.dependencies.length} dependencies</p>
              </div>
            </div>
          </motion.div>
        ))}
      </motion.div>

      {activeCategory && (
        <motion.div
          className="rounded-lg border border-border bg-card/50 p-6"
          variants={dependencyVariants}
          initial="hidden"
          animate="visible"
          exit="hidden"
        >
          <h3 className="mb-4 text-xl font-semibold">
            {categories.find((c) => c.id === activeCategory)?.name} Dependencies
          </h3>
          <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 md:grid-cols-4">
            {categories
              .find((c) => c.id === activeCategory)
              ?.dependencies.map((dep, index) => (
                <motion.div
                  key={index}
                  variants={dependencyItemVariants}
                  className="rounded-lg border border-border bg-card p-3"
                >
                  <h4 className="font-medium">{dep.name}</h4>
                  <p className="text-xs text-muted-foreground">{dep.description}</p>
                </motion.div>
              ))}
          </div>
        </motion.div>
      )}

      <div className="flex justify-center">
        <div className="rounded-md bg-muted px-3 py-1.5 text-sm text-muted-foreground">
          Click on a category to see its dependencies
        </div>
      </div>
    </div>
  );
}
