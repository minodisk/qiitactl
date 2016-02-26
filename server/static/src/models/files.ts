export type File = {
  id: string;
  name: string;
  abs: string;
  rel: string;
  children: File[];
};
