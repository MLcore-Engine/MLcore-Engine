declare module '@remix-run/router' {
  export interface RouterInit {
    basename?: string;
    // 添加其他属性
  }

  export class Router {
    readonly basename: string;
    readonly state: any;
    readonly routes: any[];
    readonly hydrationData: any;
    
    constructor(init: RouterInit);
    
    initialize(): void;
    subscribe(fn: Function): () => void;
    enableScrollRestoration(): void;
    navigate(to: string, opts?: any): void;
    fetch(key: string, routeId: string, url: string, opts?: any): void;
    revalidate(): void;
    createHref(to: string): string;
    encodeLocation(to: string): any;
    getFetcher(key: string): any;
    deleteFetcher(key: string): void;
    dispose(): void;
    _internalFetchControllers: Map<string, AbortController>;
    _internalActiveDeferreds: Map<string, any>;
  }

  export function createRouter(init: RouterInit): Router;
} 