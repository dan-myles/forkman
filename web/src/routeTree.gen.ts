/* prettier-ignore-start */

/* eslint-disable */

// @ts-nocheck

// noinspection JSUnusedGlobalSymbols

// This file is auto-generated by TanStack Router

import { createFileRoute } from '@tanstack/react-router'

// Import Routes

import { Route as rootRoute } from './routes/__root'
import { Route as IndexImport } from './routes/index'
import { Route as UnauthorizedIndexImport } from './routes/unauthorized/index'
import { Route as DashboardSidebarImport } from './routes/dashboard/_sidebar'
import { Route as DashboardSidebarOverviewIndexImport } from './routes/dashboard/_sidebar.overview/index'

// Create Virtual Routes

const DashboardImport = createFileRoute('/dashboard')()

// Create/Update Routes

const DashboardRoute = DashboardImport.update({
  path: '/dashboard',
  getParentRoute: () => rootRoute,
} as any)

const IndexRoute = IndexImport.update({
  path: '/',
  getParentRoute: () => rootRoute,
} as any)

const UnauthorizedIndexRoute = UnauthorizedIndexImport.update({
  path: '/unauthorized/',
  getParentRoute: () => rootRoute,
} as any)

const DashboardSidebarRoute = DashboardSidebarImport.update({
  id: '/_sidebar',
  getParentRoute: () => DashboardRoute,
} as any)

const DashboardSidebarOverviewIndexRoute =
  DashboardSidebarOverviewIndexImport.update({
    path: '/overview/',
    getParentRoute: () => DashboardSidebarRoute,
  } as any)

// Populate the FileRoutesByPath interface

declare module '@tanstack/react-router' {
  interface FileRoutesByPath {
    '/': {
      id: '/'
      path: '/'
      fullPath: '/'
      preLoaderRoute: typeof IndexImport
      parentRoute: typeof rootRoute
    }
    '/dashboard': {
      id: '/dashboard'
      path: '/dashboard'
      fullPath: '/dashboard'
      preLoaderRoute: typeof DashboardImport
      parentRoute: typeof rootRoute
    }
    '/dashboard/_sidebar': {
      id: '/dashboard/_sidebar'
      path: '/dashboard'
      fullPath: '/dashboard'
      preLoaderRoute: typeof DashboardSidebarImport
      parentRoute: typeof DashboardRoute
    }
    '/unauthorized/': {
      id: '/unauthorized/'
      path: '/unauthorized'
      fullPath: '/unauthorized'
      preLoaderRoute: typeof UnauthorizedIndexImport
      parentRoute: typeof rootRoute
    }
    '/dashboard/_sidebar/overview/': {
      id: '/dashboard/_sidebar/overview/'
      path: '/overview'
      fullPath: '/dashboard/overview'
      preLoaderRoute: typeof DashboardSidebarOverviewIndexImport
      parentRoute: typeof DashboardSidebarImport
    }
  }
}

// Create and export the route tree

interface DashboardSidebarRouteChildren {
  DashboardSidebarOverviewIndexRoute: typeof DashboardSidebarOverviewIndexRoute
}

const DashboardSidebarRouteChildren: DashboardSidebarRouteChildren = {
  DashboardSidebarOverviewIndexRoute: DashboardSidebarOverviewIndexRoute,
}

const DashboardSidebarRouteWithChildren =
  DashboardSidebarRoute._addFileChildren(DashboardSidebarRouteChildren)

interface DashboardRouteChildren {
  DashboardSidebarRoute: typeof DashboardSidebarRouteWithChildren
}

const DashboardRouteChildren: DashboardRouteChildren = {
  DashboardSidebarRoute: DashboardSidebarRouteWithChildren,
}

const DashboardRouteWithChildren = DashboardRoute._addFileChildren(
  DashboardRouteChildren,
)

export interface FileRoutesByFullPath {
  '/': typeof IndexRoute
  '/dashboard': typeof DashboardSidebarRouteWithChildren
  '/unauthorized': typeof UnauthorizedIndexRoute
  '/dashboard/overview': typeof DashboardSidebarOverviewIndexRoute
}

export interface FileRoutesByTo {
  '/': typeof IndexRoute
  '/dashboard': typeof DashboardSidebarRouteWithChildren
  '/unauthorized': typeof UnauthorizedIndexRoute
  '/dashboard/overview': typeof DashboardSidebarOverviewIndexRoute
}

export interface FileRoutesById {
  __root__: typeof rootRoute
  '/': typeof IndexRoute
  '/dashboard': typeof DashboardRouteWithChildren
  '/dashboard/_sidebar': typeof DashboardSidebarRouteWithChildren
  '/unauthorized/': typeof UnauthorizedIndexRoute
  '/dashboard/_sidebar/overview/': typeof DashboardSidebarOverviewIndexRoute
}

export interface FileRouteTypes {
  fileRoutesByFullPath: FileRoutesByFullPath
  fullPaths: '/' | '/dashboard' | '/unauthorized' | '/dashboard/overview'
  fileRoutesByTo: FileRoutesByTo
  to: '/' | '/dashboard' | '/unauthorized' | '/dashboard/overview'
  id:
    | '__root__'
    | '/'
    | '/dashboard'
    | '/dashboard/_sidebar'
    | '/unauthorized/'
    | '/dashboard/_sidebar/overview/'
  fileRoutesById: FileRoutesById
}

export interface RootRouteChildren {
  IndexRoute: typeof IndexRoute
  DashboardRoute: typeof DashboardRouteWithChildren
  UnauthorizedIndexRoute: typeof UnauthorizedIndexRoute
}

const rootRouteChildren: RootRouteChildren = {
  IndexRoute: IndexRoute,
  DashboardRoute: DashboardRouteWithChildren,
  UnauthorizedIndexRoute: UnauthorizedIndexRoute,
}

export const routeTree = rootRoute
  ._addFileChildren(rootRouteChildren)
  ._addFileTypes<FileRouteTypes>()

/* prettier-ignore-end */

/* ROUTE_MANIFEST_START
{
  "routes": {
    "__root__": {
      "filePath": "__root.tsx",
      "children": [
        "/",
        "/dashboard",
        "/unauthorized/"
      ]
    },
    "/": {
      "filePath": "index.tsx"
    },
    "/dashboard": {
      "filePath": "dashboard",
      "children": [
        "/dashboard/_sidebar"
      ]
    },
    "/dashboard/_sidebar": {
      "filePath": "dashboard/_sidebar.tsx",
      "parent": "/dashboard",
      "children": [
        "/dashboard/_sidebar/overview/"
      ]
    },
    "/unauthorized/": {
      "filePath": "unauthorized/index.tsx"
    },
    "/dashboard/_sidebar/overview/": {
      "filePath": "dashboard/_sidebar.overview/index.tsx",
      "parent": "/dashboard/_sidebar"
    }
  }
}
ROUTE_MANIFEST_END */
