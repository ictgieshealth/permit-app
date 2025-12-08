import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

// Middleware disabled - using client-side authentication with localStorage
// All auth checks are handled by useAuth hook in layouts
export function middleware(request: NextRequest) {
  // Let all requests through - auth is handled client-side
  return NextResponse.next();
}

export const config = {
  matcher: [
    /*
     * Match all request paths except for the ones starting with:
     * - api (API routes)
     * - _next/static (static files)
     * - _next/image (image optimization files)
     * - favicon.ico (favicon file)
     * - public folder files
     */
    "/((?!api|_next/static|_next/image|favicon.ico|.*\\.png|.*\\.jpg|.*\\.jpeg|.*\\.gif|.*\\.svg|.*\\.ico|.*\\.webp).*)",
  ],
};
