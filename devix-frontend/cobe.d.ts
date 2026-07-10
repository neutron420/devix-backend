import 'cobe';

declare module 'cobe' {
  export interface COBEOptions {
    onRender?: (state: Record<string, unknown>) => void;
  }
}
