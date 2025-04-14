import { NextResponse, type NextRequest } from 'next/server'

export const publicRoutes = [
  { path: '/', whenAuthenticated: 'next' },
  { path: '/login', whenAuthenticated: 'redirect' },
  { path: '/about', whenAuthenticated: 'next' },
  { path: '/contact', whenAuthenticated: 'next' },
  { path: '/register', whenAuthenticated: 'redirect' },
  { path: '/reset-password', whenAuthenticated: 'next' },
  { path: '/email-verification', whenAuthenticated: 'next' },
] as const

const REDIRECT_WHEN_NOT_AUTH_ROUTE = '/login'

export function middleware(request: NextRequest) {
  console.log('✅ Middleware is running')

  const pathname = request.nextUrl.pathname
  const authToken = request.cookies.get('auth-token')?.value

  const publicRoute = publicRoutes.find(route => pathname.startsWith(route.path))

  // Se não autenticado e rota pública, continua
  if (!authToken && publicRoute) return NextResponse.next()

  // Se não autenticado e rota privada, redireciona para login
  if (!authToken && !publicRoute) {
    return NextResponse.redirect(new URL(REDIRECT_WHEN_NOT_AUTH_ROUTE, request.url))
  }

  // Se autenticado e está em rota que deve redirecionar, redireciona para dashboard
  if (authToken && publicRoute?.whenAuthenticated === 'redirect') {
    return NextResponse.redirect(new URL('/dashboard', request.url))
  }

  return NextResponse.next()
}

export const config = {
  matcher: ['/((?!api|_next/static|_next/image|favicon.ico|sitemap|robots).*)']
}
