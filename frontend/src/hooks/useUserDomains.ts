import { useAuth } from "@/context/AuthContext";
import { Domain } from "@/types/domain";

export function useUserDomains(): {
  userDomains: Domain[];
  hasMultipleDomains: boolean;
  singleDomainId: number | null;
  canAccessDomain: (domainId: number) => boolean;
} {
  const { user } = useAuth();

  const userDomains = user?.domains || [];
  const hasMultipleDomains = userDomains.length > 1;
  const singleDomainId = userDomains.length === 1 ? userDomains[0].id : null;

  const canAccessDomain = (domainId: number): boolean => {
    if (!user?.domains || user.domains.length === 0) return true; // Admin can access all
    return user.domains.some((domain) => domain.id === domainId);
  };

  return {
    userDomains,
    hasMultipleDomains,
    singleDomainId,
    canAccessDomain,
  };
}
