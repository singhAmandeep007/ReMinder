export type TTechnology = {
  id: string;
  name: string;
  description: string;
  icon?: string;
};

export type TTechCategory = {
  name: string;
  color: string;
  icon: React.ReactNode;
  description: string;
  technologies: TTechnology[];
};

export type TTechStackData = {
  [key: string]: TTechCategory;
};
