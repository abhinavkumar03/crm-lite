"use client";

import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from "react";
import { toast } from "sonner";

import {
  listMyOrganizations,
  switchOrganization,
} from "@/features/organization/api";
import type { OrgMembership } from "@/features/organization/types";
import { invalidateMetadataCache } from "@/features/metadata/cache";
import { useAuth } from "@/context/AuthContext";

const ACTIVE_ORG_STORAGE_KEY = "active_organization_id";

type OrgContextValue = {
  orgs: OrgMembership[];
  activeOrg: OrgMembership | null;
  loading: boolean;
  switching: boolean;
  refreshOrgs: () => Promise<void>;
  switchWorkspace: (orgId: string) => Promise<void>;
};

const OrgContext = createContext<OrgContextValue | null>(null);

export function OrgProvider({ children }: { children: ReactNode }) {
  const auth = useAuth();
  const [orgs, setOrgs] = useState<OrgMembership[]>([]);
  const [loading, setLoading] = useState(true);
  const [switching, setSwitching] = useState(false);

  const refreshOrgs = useCallback(async () => {
    if (!auth.user) {
      setOrgs([]);
      setLoading(false);
      return;
    }
    try {
      const list = await listMyOrganizations();
      setOrgs(list);
      const active = list.find((o) => o.is_active);
      if (active) {
        localStorage.setItem(ACTIVE_ORG_STORAGE_KEY, active.id);
      }
    } catch {
      setOrgs([]);
    } finally {
      setLoading(false);
    }
  }, [auth.user]);

  useEffect(() => {
    setLoading(true);
    void refreshOrgs();
  }, [refreshOrgs]);

  const activeOrg = useMemo(
    () => orgs.find((o) => o.is_active) ?? orgs[0] ?? null,
    [orgs]
  );

  const switchWorkspace = useCallback(async (orgId: string) => {
    if (switching) return;
    try {
      setSwitching(true);
      await switchOrganization(orgId);
      localStorage.setItem(ACTIVE_ORG_STORAGE_KEY, orgId);
      invalidateMetadataCache();
      toast.success("Workspace switched");
      window.location.assign("/dashboard");
    } catch {
      toast.error("Could not switch workspace");
      setSwitching(false);
    }
  }, [switching]);

  const value = useMemo(
    () => ({
      orgs,
      activeOrg,
      loading,
      switching,
      refreshOrgs,
      switchWorkspace,
    }),
    [orgs, activeOrg, loading, switching, refreshOrgs, switchWorkspace]
  );

  return <OrgContext.Provider value={value}>{children}</OrgContext.Provider>;
}

export function useOrg(): OrgContextValue {
  const ctx = useContext(OrgContext);
  if (!ctx) {
    throw new Error("useOrg must be used within OrgProvider");
  }
  return ctx;
}

export function useOrgOptional(): OrgContextValue | null {
  return useContext(OrgContext);
}

export { ACTIVE_ORG_STORAGE_KEY };
