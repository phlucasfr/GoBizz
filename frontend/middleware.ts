import { NextResponse, type NextRequest } from 'next/server'

export const publicRoutes = [
  { path: '/', whenAuthenticated: 'next' },
  { path: '/login', whenAuthenticated: 'redirect' },
  { path: '/register', whenAuthenticated: 'redirect' },
  { path: '/reset-password', whenAuthenticated: 'next' },
  { path: '/email-verification', whenAuthenticated: 'next' },
] as const

const REDIRECT_WHEN_NOT_AUTH_ROUTE = '/login'

export function middleware(request: NextRequest) {
  const pathname = request.nextUrl.pathname
  const authToken = request.cookies.get('auth-token')?.value
  const publicRoute = publicRoutes.find((route) => route.path === pathname)


  if (!authToken && publicRoute) return NextResponse.next()

  if (!authToken && !publicRoute) {
    const redirectUrl = new URL(REDIRECT_WHEN_NOT_AUTH_ROUTE, request.url)
    return NextResponse.redirect(redirectUrl)
  }

  if (authToken && publicRoute?.whenAuthenticated === 'redirect') {
    const redirectUrl = new URL('/dashboard', request.url)
    return NextResponse.redirect(redirectUrl)
  }

  // Allow access to slug routes without authentication
  if (pathname.startsWith('/')) {
    return NextResponse.next()
  }


  return NextResponse.next()
}

export const config = {
  matcher: ['/((?!api|_next/static|_next/image|favicon.ico|sitemap|robots).*)']
}