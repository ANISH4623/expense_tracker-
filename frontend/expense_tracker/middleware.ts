import { NextResponse, type NextRequest } from 'next/server';

export function middleware(request: NextRequest) {
  const currentUser = request.cookies.get('first_name')?.value;

  // Check if the user is authenticated and trying to access the dashboard
  if (currentUser && request.nextUrl.pathname.startsWith('/dashboard')) {
    // User is already authenticated, allow access to the dashboard
    return NextResponse.next();
  }

  // User is not authenticated or trying to access a protected route
  if (!currentUser && request.nextUrl.pathname.startsWith('/dashboard')) {
    // Redirect to login if user tries to access dashboard without being logged in
    return NextResponse.redirect(new URL('/', request.url));
  }

  // Allow access to public routes if not starting with '/dashboard'
  return NextResponse.next();
}

export const config = {
  matcher: [ '/dashboard/:path*',],
};
