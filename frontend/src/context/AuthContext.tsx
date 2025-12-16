"use client";
import { createContext, useContext, useState, useEffect, ReactNode } from "react";
import { User, UserDomainRole } from "@/types/user";
import { Domain } from "@/types/domain";
import { Role } from "@/types/role";
import { usePathname } from "next/navigation";
import { authService } from "@/services/auth.service";
import { STORAGE_KEYS } from "@/config/api";

interface AuthContextType {
    user: User | null;
    currentDomain: Domain | null;
    currentRole: Role | null;
    userDomains: UserDomainRole[];
    loading: boolean;
    login: (username: string, password: string) => Promise<void>;
    logout: () => void;
    switchDomain: (domainId: number) => Promise<void>;
    refreshUser: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
    const [user, setUser] = useState<User | null>(null);
    const [currentDomain, setCurrentDomain] = useState<Domain | null>(null);
    const [currentRole, setCurrentRole] = useState<Role | null>(null);
    const [userDomains, setUserDomains] = useState<UserDomainRole[]>([]);
    const [loading, setLoading] = useState(true);
    const pathname = usePathname();

    useEffect(() => {
        checkAuth();
    }, []);

    const checkAuth = async () => {
        const token = localStorage.getItem(STORAGE_KEYS.TOKEN);
        if (token) {
            try {
                const userData = await authService.getProfile();
                const domain = authService.getStoredDomain();
                const role = authService.getStoredRole();
                const domains = authService.getStoredDomains();
                
                setUser(userData);
                setCurrentDomain(domain);
                setCurrentRole(role);
                setUserDomains(domains);
            } catch (err) {
                localStorage.removeItem(STORAGE_KEYS.TOKEN);
                setUser(null);
                setCurrentDomain(null);
                setCurrentRole(null);
                setUserDomains([]);
            }
        }
        setLoading(false);
    };

    const login = async (username: string, password: string) => {
        const response = await authService.login({ username, password });
        localStorage.setItem(STORAGE_KEYS.TOKEN, response.token);
        setUser(response.user);
        setCurrentDomain(response.current_domain);
        setCurrentRole(response.current_role);
        setUserDomains(response.domains);
    };

    const switchDomain = async (domainId: number) => {
        try {
            const response = await authService.switchDomain(domainId);
            setCurrentDomain(response.current_domain);
            setCurrentRole(response.current_role);
            // Refresh the page to reload all data with new domain context
            window.location.reload();
        } catch (err) {
            console.error("Failed to switch domain:", err);
            throw err;
        }
    };

    const logout = () => {
        authService.logout();
        setUser(null);
        setCurrentDomain(null);
        setCurrentRole(null);
        setUserDomains([]);
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
        <AuthContext.Provider value={{ 
            user, 
            currentDomain, 
            currentRole, 
            userDomains, 
            loading, 
            login, 
            logout, 
            switchDomain,
            refreshUser 
        }}>
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
