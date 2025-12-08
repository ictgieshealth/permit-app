"use client";

import { useEffect, useState } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import { authService } from '@/services/auth.service';

export function useAuth(requireAuth: boolean = true) {
    const router = useRouter();
    const pathname = usePathname();
    const [isChecking, setIsChecking] = useState(true);

    useEffect(() => {
        // Only run on client-side
        if (typeof window === 'undefined') return;

        const isAuthenticated = authService.isAuthenticated();

        if (requireAuth && !isAuthenticated) {
            // Not authenticated, redirect to signin
            router.replace(`/signin?from=${pathname}`);
        } else if (!requireAuth && isAuthenticated) {
            // Already authenticated, redirect to dashboard
            router.replace('/');
        } else {
            setIsChecking(false);
        }
    }, [requireAuth, router, pathname]);

    return {
        isAuthenticated: authService.isAuthenticated(),
        user: authService.getStoredUser(),
        domain: authService.getStoredDomain(),
        isChecking,
    };
}
