"use client";
import { createContext, useContext, useState, useEffect, ReactNode } from "react";
import { User } from "@/types/user";
import { usePathname } from "next/navigation";
import { authService } from "@/services/auth.service";

interface AuthContextType {
    user: User | null;
    loading: boolean;
    login: (username: string, password: string) => Promise<void>;
    logout: () => void;
    refreshUser: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
    const [user, setUser] = useState<User | null>(null);
    const [loading, setLoading] = useState(true);
    const pathname = usePathname();

    useEffect(() => {
        checkAuth();
    }, []);

    const checkAuth = async () => {
        const token = localStorage.getItem("auth_token");
        if (token) {
            try {
                const userData = await authService.getProfile();
                setUser(userData);
            } catch (err) {
                localStorage.removeItem("token");
                setUser(null);
            }
        }
        setLoading(false);
    };

    const login = async (username: string, password: string) => {
        const response = await authService.login({username, password});
        console.log('Login response:', response);
        localStorage.setItem("token", response.token);
        setUser(response.user);
    };

    const logout = () => {
        localStorage.removeItem("token");
        setUser(null);
        window.location.href = "/auth/signin";
    };

    const refreshUser = async () => {
        try {
            const userData = await authService.getProfile();
            setUser(userData);
        } catch (err) {
            console.error("Failed to refresh user:", err);
        }
    };

    return (
        <AuthContext.Provider value={{ user, loading, login, logout, refreshUser }}>
            {children}
        </AuthContext.Provider>
    );
}

export function useAuth() {
    const context = useContext(AuthContext);
    if (context === undefined) {
        throw new Error("useAuth must be used within an AuthProvider");
    }
    return context;
}
