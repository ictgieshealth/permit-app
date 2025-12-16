"use client";
import React, { useEffect, useRef, useState, useCallback } from "react";
import Link from "next/link";
import Image from "next/image";
import { usePathname } from "next/navigation";
import { useSidebar } from "../context/SidebarContext";
import { useAuth } from "../context/AuthContext";
import {
  BoxCubeIcon,
  CalenderIcon,
  ChevronDownIcon,
  GridIcon,
  HorizontaLDots,
  ListIcon,
  PageIcon,
  PieChartIcon,
  PlugInIcon,
  TableIcon,
  UserCircleIcon,
} from "../icons/index";
import SidebarWidget from "./SidebarWidget";
import { Menu } from "@/types/menu";
import { menuService } from "@/services/menu.service";

type NavItem = {
  name: string;
  icon: React.ReactNode;
  path?: string;
  subItems?: { name: string; path: string; pro?: boolean; new?: boolean }[];
};

const iconMap: Record<string, React.ReactNode> = {
  GridIcon: <GridIcon />,
  ListIcon: <ListIcon />,
  PageIcon: <PageIcon />,
  BoxCubeIcon: <BoxCubeIcon />,
  UserCircleIcon: <UserCircleIcon />,
  CalenderIcon: <CalenderIcon />,
  PieChartIcon: <PieChartIcon />,
  PlugInIcon: <PlugInIcon />,
  TableIcon: <TableIcon />,
};

const AppSidebar: React.FC = () => {
  const { isExpanded, isMobileOpen, isHovered, setIsHovered } = useSidebar();
  const { user, currentDomain, userDomains, switchDomain, loading: authLoading } = useAuth();
  const pathname = usePathname();
  const [navItems, setNavItems] = useState<NavItem[]>([]);
  const [ticketingItems, setTicketingItems] = useState<NavItem[]>([]);
  const [othersItems, setOthersItems] = useState<NavItem[]>([]);
  const [menuLoading, setMenuLoading] = useState(true);
  const [switchingDomain, setSwitchingDomain] = useState(false);

  useEffect(() => {
    if (!authLoading && user) {
      loadUserMenus();
    }
  }, [authLoading, user]);

  const loadUserMenus = async () => {
    try {
      setMenuLoading(true);
      const menus = await menuService.getUserMenus();
      
      // Separate main menu items from master data and ticketing items
      const masterDataPaths = ['/domains', '/divisions', '/users', '/menus', '/permit-types', '/roles'];
      const ticketingPaths = ['/tasks', '/tasks/task-requests', '/task-requests', '/projects'];
      
      const mainMenus = menus.filter((menu) => !menu.parent_id && !masterDataPaths.includes(menu.path) && !ticketingPaths.includes(menu.path));
      const masterDataMenus = menus.filter((menu) => !menu.parent_id && masterDataPaths.includes(menu.path));
      const ticketingMenus = menus.filter((menu) => !menu.parent_id && ticketingPaths.includes(menu.path));
      
      const convertMenuToNavItem = (menu: Menu): NavItem => ({
        name: menu.name,
        icon: menu.icon ? iconMap[menu.icon] || <PageIcon /> : <PageIcon />,
        path: menu.path,
        subItems: menu.children?.map((child) => ({
          name: child.name,
          path: child.path,
        })),
      });

      setNavItems(mainMenus.map(convertMenuToNavItem));
      setTicketingItems(ticketingMenus.map(convertMenuToNavItem));
      setOthersItems(masterDataMenus.map(convertMenuToNavItem));
    } catch (err) {
      console.error("Failed to load user menus:", err);
      // Fallback to default menus if loading fails
      setNavItems([
        { icon: <GridIcon />, name: "Dashboard", path: "/" },
        { icon: <ListIcon />, name: "Permits", path: "/permits" },
      ]);
      setTicketingItems([
        { icon: <PageIcon />, name: "Tasks", path: "/tasks" },
        { icon: <PageIcon />, name: "Task Requests", path: "/task-requests" },
        { icon: <PageIcon />, name: "Projects", path: "/projects" },
      ]);
      setOthersItems([]);
    } finally {
      setMenuLoading(false);
    }
  };

  const renderMenuItems = (
    navItems: NavItem[],
    menuType: "main" | "ticketing" | "others"
  ) => (
    <ul className="flex flex-col gap-4">
      {navItems.map((nav, index) => (
        <li key={nav.name}>
          {nav.subItems ? (
            <button
              onClick={() => handleSubmenuToggle(index, menuType)}
              className={`menu-item group  ${openSubmenu?.type === menuType && openSubmenu?.index === index
                ? "menu-item-active"
                : "menu-item-inactive"
                } cursor-pointer ${!isExpanded && !isHovered
                  ? "lg:justify-center"
                  : "lg:justify-start"
                }`}
            >
              <span
                className={` ${openSubmenu?.type === menuType && openSubmenu?.index === index
                  ? "menu-item-icon-active"
                  : "menu-item-icon-inactive"
                  }`}
              >
                {nav.icon}
              </span>
              {(isExpanded || isHovered || isMobileOpen) && (
                <span className={`menu-item-text`}>{nav.name}</span>
              )}
              {(isExpanded || isHovered || isMobileOpen) && (
                <ChevronDownIcon
                  className={`ml-auto w-5 h-5 transition-transform duration-200  ${openSubmenu?.type === menuType &&
                    openSubmenu?.index === index
                    ? "rotate-180 text-brand-500"
                    : ""
                    }`}
                />
              )}
            </button>
          ) : (
            nav.path && (
              <Link
                href={nav.path}
                className={`menu-item group ${isActive(nav.path) ? "menu-item-active" : "menu-item-inactive"
                  }`}
              >
                <span
                  className={`${isActive(nav.path)
                    ? "menu-item-icon-active"
                    : "menu-item-icon-inactive"
                    }`}
                >
                  {nav.icon}
                </span>
                {(isExpanded || isHovered || isMobileOpen) && (
                  <span className={`menu-item-text`}>{nav.name}</span>
                )}
              </Link>
            )
          )}
          {nav.subItems && (isExpanded || isHovered || isMobileOpen) && (
            <div
              ref={(el) => {
                subMenuRefs.current[`${menuType}-${index}`] = el;
              }}
              className="overflow-hidden transition-all duration-300"
              style={{
                height:
                  openSubmenu?.type === menuType && openSubmenu?.index === index
                    ? `${subMenuHeight[`${menuType}-${index}`]}px`
                    : "0px",
              }}
            >
              <ul className="mt-2 space-y-1 ml-9">
                {nav.subItems.map((subItem) => (
                  <li key={subItem.name}>
                    <Link
                      href={subItem.path}
                      className={`menu-dropdown-item ${isActive(subItem.path)
                        ? "menu-dropdown-item-active"
                        : "menu-dropdown-item-inactive"
                        }`}
                    >
                      {subItem.name}
                      <span className="flex items-center gap-1 ml-auto">
                        {subItem.new && (
                          <span
                            className={`ml-auto ${isActive(subItem.path)
                              ? "menu-dropdown-badge-active"
                              : "menu-dropdown-badge-inactive"
                              } menu-dropdown-badge `}
                          >
                            new
                          </span>
                        )}
                        {subItem.pro && (
                          <span
                            className={`ml-auto ${isActive(subItem.path)
                              ? "menu-dropdown-badge-active"
                              : "menu-dropdown-badge-inactive"
                              } menu-dropdown-badge `}
                          >
                            pro
                          </span>
                        )}
                      </span>
                    </Link>
                  </li>
                ))}
              </ul>
            </div>
          )}
        </li>
      ))}
    </ul>
  );

  const [openSubmenu, setOpenSubmenu] = useState<{
    type: "main" | "ticketing" | "others";
    index: number;
  } | null>(null);
  const [subMenuHeight, setSubMenuHeight] = useState<Record<string, number>>(
    {}
  );
  const subMenuRefs = useRef<Record<string, HTMLDivElement | null>>({});

  // const isActive = (path: string) => path === pathname;
  const isActive = useCallback((path: string) => path === pathname, [pathname]);

  useEffect(() => {
    // Check if the current path matches any submenu item
    let submenuMatched = false;
    ["main", "ticketing", "others"].forEach((menuType) => {
      const items = menuType === "main" ? navItems : menuType === "ticketing" ? ticketingItems : othersItems;
      items.forEach((nav, index) => {
        if (nav.subItems) {
          nav.subItems.forEach((subItem) => {
            if (isActive(subItem.path)) {
              setOpenSubmenu({
                type: menuType as "main" | "ticketing" | "others",
                index,
              });
              submenuMatched = true;
            }
          });
        }
      });
    });

    // If no submenu item matches, close the open submenu
    if (!submenuMatched) {
      setOpenSubmenu(null);
    }
  }, [pathname, isActive, navItems, ticketingItems, othersItems]);

  useEffect(() => {
    // Set the height of the submenu items when the submenu is opened
    if (openSubmenu !== null) {
      const key = `${openSubmenu.type}-${openSubmenu.index}`;
      if (subMenuRefs.current[key]) {
        setSubMenuHeight((prevHeights) => ({
          ...prevHeights,
          [key]: subMenuRefs.current[key]?.scrollHeight || 0,
        }));
      }
    }
  }, [openSubmenu]);

  const handleSubmenuToggle = (index: number, menuType: "main" | "ticketing" | "others") => {
    setOpenSubmenu((prevOpenSubmenu) => {
      if (
        prevOpenSubmenu &&
        prevOpenSubmenu.type === menuType &&
        prevOpenSubmenu.index === index
      ) {
        return null;
      }
      return { type: menuType, index };
    });
  };

  const handleDomainSwitch = async (domainId: number) => {
    if (domainId === currentDomain?.id || switchingDomain) return;
    
    try {
      setSwitchingDomain(true);
      await switchDomain(domainId);
    } catch (err) {
      console.error("Failed to switch domain:", err);
      alert("Failed to switch domain. Please try again.");
    } finally {
      setSwitchingDomain(false);
    }
  };

  if (menuLoading || authLoading) {
    return (
      <aside className="fixed mt-16 flex flex-col lg:mt-0 top-0 px-5 left-0 bg-white dark:bg-gray-900 dark:border-gray-800 text-gray-900 h-screen transition-all duration-300 ease-in-out z-50 border-r border-gray-200 w-[290px] lg:translate-x-0 -translate-x-full">
        <div className="flex items-center justify-center h-full">
          <div className="text-gray-500">Loading...</div>
        </div>
      </aside>
    );
  }

  return (
    <aside
      className={`fixed mt-16 flex flex-col lg:mt-0 top-0 px-5 left-0 bg-white dark:bg-gray-900 dark:border-gray-800 text-gray-900 h-screen transition-all duration-300 ease-in-out z-50 border-r border-gray-200 
        ${isExpanded || isMobileOpen
          ? "w-[290px]"
          : isHovered
            ? "w-[290px]"
            : "w-[90px]"
        }
        ${isMobileOpen ? "translate-x-0" : "-translate-x-full"}
        lg:translate-x-0`}
      onMouseEnter={() => !isExpanded && setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
    >
      <div
        className={`py-8 flex  ${!isExpanded && !isHovered ? "lg:justify-center" : "justify-start"
          }`}
      >
        <Link href="/">
          {isExpanded || isHovered || isMobileOpen ? (
            <>
              <Image
                className="dark:hidden"
                src="/images/logo/logo.png"
                alt="Logo"
                width={240}
                height={64}
              />
              <Image
                className="hidden dark:block"
                src="/images/logo/logo-dark.svg"
                alt="Logo"
                width={150}
                height={40}
              />
            </>
          ) : (
            <Image
              src="/images/logo/logo-icon.svg"
              alt="Logo"
              width={32}
              height={32}
            />
          )}
        </Link>
      </div>

      {/* Domain Switcher */}
      {userDomains && userDomains.length > 1 && (isExpanded || isHovered || isMobileOpen) && (
        <div className="mb-6 px-2">
          <div className="relative">
            <label className="block text-xs font-medium text-gray-500 dark:text-gray-400 mb-2 px-1">
              Active Domain
            </label>
            <div className="relative">
              <select
                value={currentDomain?.id || ""}
                onChange={(e) => handleDomainSwitch(Number(e.target.value))}
                disabled={switchingDomain}
                className="w-full pl-10 pr-8 py-2.5 text-sm font-medium border-0 rounded-xl bg-gradient-to-r from-gray-50 to-gray-100 dark:from-gray-800 dark:to-gray-700 text-gray-900 dark:text-white shadow-sm hover:shadow-md focus:outline-none focus:ring-2 focus:ring-brand-500 focus:ring-offset-2 dark:focus:ring-offset-gray-900 disabled:opacity-60 disabled:cursor-not-allowed transition-all duration-200 appearance-none cursor-pointer"
              >
                {userDomains.map((udr) => (
                  <option key={udr.domain_id} value={udr.domain_id}>
                    {udr.domain?.name || `Domain ${udr.domain_id}`}
                  </option>
                ))}
              </select>
              <div className="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none">
                <svg className="w-5 h-5 text-gray-500 dark:text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
                </svg>
              </div>
              <div className="absolute inset-y-0 right-0 flex items-center pr-3 pointer-events-none">
                <svg className={`w-4 h-4 text-gray-500 dark:text-gray-400 transition-transform ${switchingDomain ? 'animate-spin' : ''}`} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  {switchingDomain ? (
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                  ) : (
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 9l4-4 4 4m0 6l-4 4-4-4" />
                  )}
                </svg>
              </div>
            </div>
            {switchingDomain && (
              <p className="mt-2 text-xs text-center text-brand-600 dark:text-brand-400 font-medium animate-pulse">
                Switching domain...
              </p>
            )}
          </div>
        </div>
      )}
      
      <div className="flex flex-col overflow-y-auto duration-300 ease-linear no-scrollbar">
        <nav className="mb-6">
          <div className="flex flex-col gap-4">
            <div>
              <h2
                className={`mb-4 text-xs uppercase flex leading-[20px] text-gray-400 ${!isExpanded && !isHovered
                  ? "lg:justify-center"
                  : "justify-start"
                  }`}
              >
                {isExpanded || isHovered || isMobileOpen ? (
                  "Menu"
                ) : (
                  <HorizontaLDots />
                )}
              </h2>
              {renderMenuItems(navItems, "main")}
            </div>

            <div className="">
              <h2
                className={`mb-4 text-xs uppercase flex leading-[20px] text-gray-400 ${!isExpanded && !isHovered
                  ? "lg:justify-center"
                  : "justify-start"
                  }`}
              >
                {isExpanded || isHovered || isMobileOpen ? (
                  "Ticketing"
                ) : (
                  <HorizontaLDots />
                )}
              </h2>
              {ticketingItems.length > 0 ? (
                renderMenuItems(ticketingItems, "ticketing")
              ) : (
                <div className="text-xs text-gray-500">No ticketing items</div>
              )}
            </div>

            <div className="">
              <h2
                className={`mb-4 text-xs uppercase flex leading-[20px] text-gray-400 ${!isExpanded && !isHovered
                  ? "lg:justify-center"
                  : "justify-start"
                  }`}
              >
                {isExpanded || isHovered || isMobileOpen ? (
                  "Master Data"
                ) : (
                  <HorizontaLDots />
                )}
              </h2>
              {othersItems.length > 0 ? (
                renderMenuItems(othersItems, "others")
              ) : (
                <div className="text-xs text-gray-500">No menu items</div>
              )}
            </div>
          </div>
        </nav>
      </div>
    </aside>
  );
};

export default AppSidebar;
