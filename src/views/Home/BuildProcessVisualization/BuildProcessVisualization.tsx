import { useState, useEffect, useRef } from "react";
import { motion } from "framer-motion";
import { Code, FileCode, Package, Archive, Loader2, CheckCircle, Cpu, Globe } from "lucide-react";

import { Button } from "components";

export function BuildProcessVisualization() {
  const [currentStep, setCurrentStep] = useState(0);
  const [isPlaying, setIsPlaying] = useState(false);
  const [isComplete, setIsComplete] = useState(false);
  const intervalRef = useRef<NodeJS.Timeout | null>(null);

  const steps = [
    {
      id: "source",
      title: "Source Code",
      icon: <FileCode className="h-8 w-8 text-primary" />,
      description: "TypeScript, React, and CSS files that make up the application",
    },
    {
      id: "transpile",
      title: "Transpilation",
      icon: <Code className="h-8 w-8 text-primary" />,
      description: "Converting TypeScript to JavaScript and JSX to React elements",
    },
    {
      id: "bundle",
      title: "Bundling",
      icon: <Package className="h-8 w-8 text-primary" />,
      description: "Combining modules and dependencies into optimized bundles",
    },
    {
      id: "optimize",
      title: "Optimization",
      icon: <Cpu className="h-8 w-8 text-primary" />,
      description: "Minifying, tree-shaking, and code-splitting for performance",
    },
    {
      id: "package",
      title: "Packaging",
      icon: <Archive className="h-8 w-8 text-primary" />,
      description: "Creating a deployable package with all assets",
    },
    {
      id: "deploy",
      title: "Deployment",
      icon: <Globe className="h-8 w-8 text-primary" />,
      description: "Deploying to production servers and CDN",
    },
  ];

  // Control animation playback
  useEffect(() => {
    if (isPlaying) {
      intervalRef.current = setInterval(() => {
        setCurrentStep((prev) => {
          if (prev >= steps.length - 1) {
            setIsPlaying(false);
            setIsComplete(true);
            return prev;
          }
          return prev + 1;
        });
      }, 2000);
    } else if (intervalRef.current) {
      clearInterval(intervalRef.current);
    }

    return () => {
      if (intervalRef.current) {
        clearInterval(intervalRef.current);
      }
    };
  }, [isPlaying, steps.length]);

  const handlePlay = () => {
    if (isComplete) {
      setCurrentStep(0);
      setIsComplete(false);
    }
    setIsPlaying(true);
  };

  const handlePause = () => {
    setIsPlaying(false);
  };

  const handleReset = () => {
    setCurrentStep(0);
    setIsPlaying(false);
    setIsComplete(false);
  };

  return (
    <div className="flex flex-col space-y-8">
      {/* Progress bar */}
      <div className="relative h-2 w-full rounded-full bg-muted">
        <motion.div
          className="absolute left-0 top-0 h-full rounded-full bg-primary"
          initial={{ width: "0%" }}
          animate={{ width: `${(currentStep / steps.length - 1) * 100}%` }}
          transition={{ type: "spring", stiffness: 100, damping: 20 }}
        />
      </div>

      {/* Build process visualization */}
      <div className="relative mx-auto flex w-full max-w-3xl flex-col items-center justify-center">
        {/* Connection lines */}
        <div className="absolute left-1/2 top-[53px] h-[89%] w-0.5 -translate-x-1/2 bg-border" />

        {/* Steps */}
        <div className="relative z-10 flex w-full flex-col space-y-12">
          {steps.map((step, index) => (
            <div
              key={step.id}
              className={`flex ${index % 2 === 0 ? "flex-row" : "flex-row-reverse"} items-center justify-center gap-4`}
            >
              {/* Step content */}
              <motion.div
                className={`w-[calc(50%-2rem)] rounded-lg border border-border bg-card p-4 shadow-md ${
                  currentStep >= index ? "border-primary/50" : ""
                }`}
                initial={{ opacity: 0, y: 20 }}
                animate={{
                  opacity: currentStep >= index ? 1 : 0.5,
                  y: 0,
                  scale: currentStep === index ? 1.05 : 1,
                }}
                transition={{
                  type: "spring",
                  stiffness: 300,
                  damping: 30,
                  delay: index * 0.1,
                }}
              >
                <h3 className="mb-2 text-lg font-semibold">{step.title}</h3>
                <p className="text-sm text-muted-foreground">{step.description}</p>
              </motion.div>

              {/* Center node */}
              <motion.div
                className={`relative z-20 flex h-12 w-12 items-center justify-center rounded-full border-2 bg-card ${
                  currentStep >= index ? "border-primary" : "border-border"
                }`}
                initial={{ scale: 0 }}
                animate={{
                  scale: 1,
                }}
                transition={{
                  type: "spring",
                  stiffness: 300,
                  damping: 20,
                  delay: index * 0.1,
                }}
              >
                {currentStep > index || (currentStep === steps.length - 1 && isComplete) ? (
                  <CheckCircle className="h-6 w-6 text-primary" />
                ) : currentStep === index ? (
                  <motion.div
                    animate={{ rotate: isPlaying ? 360 : 0 }}
                    transition={{ duration: 2, repeat: isPlaying ? Number.POSITIVE_INFINITY : 0, ease: "linear" }}
                  >
                    <Loader2 className="h-6 w-6 text-primary" />
                  </motion.div>
                ) : (
                  step.icon
                )}

                {/* Step number */}
                <div className="absolute -right-2 -top-2 flex h-6 w-6 items-center justify-center rounded-full bg-muted text-xs font-semibold">
                  {index + 1}
                </div>
              </motion.div>

              {/* Empty div for layout */}
              <div className="w-[calc(50%-2rem)]" />
            </div>
          ))}
        </div>
      </div>

      {/* Controls */}
      <div className="flex justify-center space-x-4">
        {isPlaying ? (
          <Button
            onClick={handlePause}
            variant="outline"
          >
            Pause
          </Button>
        ) : (
          <Button onClick={handlePlay}>{isComplete ? "Replay" : "Start Build Process"}</Button>
        )}
        <Button
          onClick={handleReset}
          variant="outline"
          disabled={currentStep === 0}
        >
          Reset
        </Button>
      </div>

      {/* Status */}
      <div className="text-center text-sm text-muted-foreground">
        {isComplete ? (
          <span className="flex items-center justify-center text-primary">
            <CheckCircle className="mr-2 h-4 w-4" />
            Build process completed successfully!
          </span>
        ) : isPlaying ? (
          <span className="flex items-center justify-center">
            <Loader2 className="mr-2 h-4 w-4 animate-spin" />
            Building: {steps[currentStep].title}...
          </span>
        ) : (
          <span>Click "Start Build Process" to begin the animation</span>
        )}
      </div>
    </div>
  );
}
