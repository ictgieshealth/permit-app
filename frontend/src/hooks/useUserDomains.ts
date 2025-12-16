import { useAuth } from "@/context/AuthContext";
import { Domain } from "@/types/domain";

export function useUserDomains(): {
  userDomains: Domain[];
  hasMultipleDomains: boolean;
  singleDomainId: number | null;
  canAccessDomain: (domainId: number) => boolean;
} {
  const { user } = useAuth();

  // Extract unique domains from domain_roles
  const domainMap = new Map<number, Domain>();
  user?.domain_roles?.forEach((dr) => {
    if (dr.domain && !domainMap.has(dr.domain.id)) {
      domainMap.set(dr.domain.id, dr.domain);
    }
  });
  
  const userDomains = Array.from(domainMap.values());
  const hasMultipleDomains = userDomains.length > 1;
  const singleDomainId = userDomains.length === 1 ? userDomains[0].id : null;

  const canAccessDomain = (domainId: number): boolean => {
    if (!user?.domain_roles || user.domain_roles.length === 0) return true; // Admin can access all
    return user.domain_roles.some((dr) => dr.domain_id === domainId);
  };

  return {
    userDomains,
    hasMultipleDomains,
    singleDomainId,
    canAccessDomain,
  };
}
