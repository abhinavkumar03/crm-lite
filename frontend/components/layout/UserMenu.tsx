"use client";

import {
  useEffect,
  useRef,
  useState,
} from "react";

import Link from "next/link";

import { useRouter } from "next/navigation";

import {
  ChevronDown,
  LogOut,
  LayoutDashboard,
  Compass,
  Map,
  CircleHelp,
  Building2,
  Check,
} from "lucide-react";
import { toast } from "sonner";

import { useAuth } from "@/context/AuthContext";
import { useDemo } from "@/features/demo/DemoProvider";
import { useTour } from "@/features/tour/TourProvider";
import {
  listMyOrganizations,
  switchOrganization,
} from "@/features/organization/api";
import { OrgMembership } from "@/features/organization/types";

type Props = {
  mobile?: boolean;
  showDashboard?: boolean;
};

export default function UserMenu({
  mobile = false,
  showDashboard = false,
}: Props)  {
  const auth = useAuth();

  const tour = useTour();
  const demo = useDemo();

  const router = useRouter();

  const [open, setOpen] =
    useState(false);

  const [orgs, setOrgs] = useState<OrgMembership[]>([]);
  const [switching, setSwitching] = useState(false);

  const ref =
    useRef<HTMLDivElement>(null);

  function logout() {
    auth.logout();

    router.replace("/login");
  }

  useEffect(() => {
    if (!auth.user) return;
    let active = true;
    (async () => {
      try {
        const list = await listMyOrganizations();
        if (active) setOrgs(list);
      } catch {
        // Single-org demos can ignore switcher failures.
      }
    })();
    return () => {
      active = false;
    };
  }, [auth.user]);

  async function handleSwitchOrg(orgId: string) {
    if (switching) return;
    try {
      setSwitching(true);
      await switchOrganization(orgId);
      setOrgs((prev) =>
        prev.map((o) => ({ ...o, is_active: o.id === orgId }))
      );
      toast.success("Organization switched");
      setOpen(false);
      router.refresh();
      window.location.reload();
    } catch {
      toast.error("Could not switch organization");
    } finally {
      setSwitching(false);
    }
  }

  useEffect(() => {
    function handleClickOutside(
      event: MouseEvent
    ) {
      if (
        ref.current &&
        !ref.current.contains(
          event.target as Node
        )
      ) {
        setOpen(false);
      }
    }

    function handleEscape(
      event: KeyboardEvent
    ) {
      if (event.key === "Escape") {
        setOpen(false);
      }
    }

    document.addEventListener(
      "mousedown",
      handleClickOutside
    );

    window.addEventListener(
      "keydown",
      handleEscape
    );

    return () => {
      document.removeEventListener(
        "mousedown",
        handleClickOutside
      );

      window.removeEventListener(
        "keydown",
        handleEscape
      );
    };
  }, []);

  const initials =
    auth.user?.name
      ?.split(" ")
      .map((n) => n[0])
      .join("")
      .substring(0, 2)
      .toUpperCase() ?? "U";

      if (mobile) {
  return (
    <button
      onClick={logout}
      className="
      flex
      w-full
      items-center
      justify-center
      rounded-2xl
      border
      border-red-200
      px-5
      py-3
      font-medium
      text-red-500
      transition
      hover:bg-red-50
      "
    >
      <LogOut
        size={18}
        className="mr-2"
      />

      Logout
    </button>
  );
}

  return (
    <div
      ref={ref}
      className="relative"
    >
      {/* Trigger */}

      <button
        onClick={() =>
          setOpen((prev) => !prev)
        }
        className="
        flex
        items-center
        gap-3
        rounded-2xl 
        bg-white
        px-3
        py-2
        transition-all
        duration-200

        hover:border-emerald-300
        hover:bg-slate-50
        hover:shadow-md
        "
      >
        <div
          className="
          flex
          h-10
          w-10
          items-center
          justify-center
          rounded-full
          bg-gradient-to-br
          from-emerald-500
          to-teal-500
          font-semibold
          text-white
          "
        >
          {initials}
        </div>

        {/* <div className="hidden text-left lg:block">
          <h4 className="max-w-[160px] truncate text-sm font-semibold text-slate-900">
            {auth.user?.name}
          </h4>

          <p className="max-w-[160px] truncate text-xs text-slate-500">
            {auth.user?.email}
          </p>
        </div> */}

        <ChevronDown
          size={18}
          className={`text-slate-400 transition-transform duration-200 ${
            open ? "rotate-180" : ""
          }`}
        />
      </button>

      {/* Dropdown */}

      {open && (
        <div
          className="
          absolute
          right-0
          mt-3
          w-72
          overflow-hidden
          rounded-3xl
          border
          border-slate-200
          bg-white
          shadow-2xl
          "
        >
          {/* Profile */}

          <div className="border-b border-slate-100 px-6 py-6 text-center">
            <div
              className="
              mx-auto
              flex
              h-16
              w-16
              items-center
              justify-center
              rounded-full
              bg-gradient-to-br
              from-emerald-500
              to-teal-500
              text-xl
              font-bold
              text-white
              "
            >
              {initials}
            </div>

            <h3 className="mt-4 truncate text-lg font-semibold text-slate-900">
              {auth.user?.name}
            </h3>

            <p className="mt-1 truncate text-sm text-slate-500">
              {auth.user?.email}
            </p>
          </div>

          {/* Actions */}

          <div className="p-2">

    {orgs.length === 0 && (
      <div className="mb-1 border-b border-slate-100 px-2 pb-2">
        <Link
          href="/onboarding/organization"
          onClick={() => setOpen(false)}
          className="flex w-full items-center gap-3 rounded-2xl px-4 py-2.5 text-left text-emerald-700 transition hover:bg-emerald-50"
        >
          <Building2 size={18} />
          <span className="text-sm font-semibold">Create workspace</span>
        </Link>
      </div>
    )}

    {orgs.length > 1 && (
      <div className="mb-1 border-b border-slate-100 px-2 pb-2">
        <p className="px-2 py-1 text-xs font-semibold uppercase tracking-wide text-slate-400">
          Organization
        </p>
        {orgs.map((org) => (
          <button
            key={org.id}
            type="button"
            disabled={switching || org.is_active}
            onClick={() => handleSwitchOrg(org.id)}
            className="flex w-full items-center gap-3 rounded-2xl px-4 py-2.5 text-left text-slate-700 transition hover:bg-slate-100 disabled:opacity-60"
          >
            <Building2 size={18} />
            <span className="min-w-0 flex-1 truncate text-sm">{org.name}</span>
            {org.is_active ? <Check size={16} className="text-emerald-500" /> : null}
          </button>
        ))}
      </div>
    )}

    {showDashboard && (
        <Link
            href="/dashboard"
            onClick={() =>
                setOpen(false)
            }
            className="
            flex
            items-center
            gap-3
            rounded-2xl
            px-4
            py-3
            text-slate-700
            transition
            hover:bg-slate-100
            "
        >
            <LayoutDashboard
                size={18}
            />

            Dashboard
        </Link>
    )}

    {demo && (
        <button
            onClick={() => {
                setOpen(false);
                demo.openLauncher();
            }}
            className="
            flex
            w-full
            items-center
            gap-3
            rounded-2xl
            px-4
            py-3
            text-slate-700
            transition
            hover:bg-slate-100
            "
        >
            <Compass size={18} />

            Explore CRM
        </button>
    )}

    {tour && (
        <button
            onClick={() => {
                setOpen(false);
                tour.restart();
            }}
            className="
            flex
            w-full
            items-center
            gap-3
            rounded-2xl
            px-4
            py-3
            text-slate-700
            transition
            hover:bg-slate-100
            "
        >
            <Map size={18} />

            Take a tour
        </button>
    )}

    <Link
      href="/help"
      onClick={() => setOpen(false)}
      className="
      flex
      items-center
      gap-3
      rounded-2xl
      px-4
      py-3
      text-slate-700
      transition
      hover:bg-slate-100
      "
    >
      <CircleHelp size={18} />
      How it works
    </Link>

    <button
        onClick={logout}
        className="
        flex
        w-full
        items-center
        gap-3
        rounded-2xl
        px-4
        py-3
        text-red-600
        transition
        hover:bg-red-50
        "
    >
        <LogOut size={18} />

        Logout
    </button>

</div>
        </div>
      )}
    </div>
  );
}