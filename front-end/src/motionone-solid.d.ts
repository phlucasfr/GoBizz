declare module '@motionone/solid' {
    import { JSX } from 'solid-js'
  
    export interface MotionProps {
      initial?: Record<string, any>
      animate?: Record<string, any>
      whileTap?: Record<string, any>
      transition?: Record<string, any>
      whileHover?: Record<string, any>
    }
  
    export const Motion: {
      div: (props: MotionProps & JSX.HTMLAttributes<HTMLDivElement>) => JSX.Element
      button: (props: MotionProps & JSX.HTMLAttributes<HTMLButtonElement>) => JSX.Element
    }
  }
  
  