"use client";
import { useState, useEffect } from "react";
import { useParams } from "next/navigation";
import Button from "@/components/ui/button/Button";
import Link from "next/link";
import { Menu } from "@/types/menu";
import { menuService } from "@/services/menu.service";

export default function ViewMenuPage() {
  const params = useParams();
  const id = parseInt(params.id as string);
  const [menu, setMenu] = useState<Menu | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    loadMenu();
  }, []);

  const loadMenu = async () => {
    try {
      setLoading(true);
      const data = await menuService.getById(id);
      setMenu(data);
    } catch (err: any) {
      setError(err.message || "Failed to load menu");
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="max-w-2xl p-6 mx-auto">
        <div className="text-center text-gray-500">Loading...</div>
      </div>
    );
  }

  if (error || !menu) {
    return (
      <div className="max-w-2xl p-6 mx-auto">
        <div className="p-4 text-sm border rounded-lg bg-error-50 border-error-200 text-error-700">
          {error || "Menu not found"}
        </div>
        <div className="mt-4">
          <Link href="/menus">
            <Button>Back to Menus</Button>
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-2xl p-6 mx-auto">
      <div className="mb-6">
        <div className="flex items-center gap-2 mb-2">
          <Link href="/menus" className="text-gray-500 hover:text-gray-700 dark:text-gray-400">
            ← Back
          </Link>
        </div>
        <h1 className="text-2xl font-bold text-gray-800 dark:text-white">Menu Details</h1>
        <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
          View menu information and role assignments
        </p>
      </div>

      <div className="p-6 space-y-6 bg-white border border-gray-200 rounded-lg dark:bg-gray-dark dark:border-gray-800">
        <div>
          <div className="text-sm text-gray-500 dark:text-gray-400">Menu Name</div>
          <div className="mt-1 text-lg font-medium text-gray-800 dark:text-white">
            {menu.name}
          </div>
        </div>

        <div>
          <div className="text-sm text-gray-500 dark:text-gray-400">Path</div>
          <div className="mt-1 text-gray-800 dark:text-white">
            <code className="px-2 py-1 bg-gray-100 dark:bg-gray-800 rounded">{menu.path}</code>
          </div>
        </div>

        {menu.icon && (
          <div>
            <div className="text-sm text-gray-500 dark:text-gray-400">Icon</div>
            <div className="mt-1 text-gray-800 dark:text-white">
              <code className="px-2 py-1 bg-gray-100 dark:bg-gray-800 rounded">{menu.icon}</code>
            </div>
          </div>
        )}

        <div>
          <div className="text-sm text-gray-500 dark:text-gray-400">Order Index</div>
          <div className="mt-1 text-gray-800 dark:text-white">{menu.order_index}</div>
        </div>

        <div>
          <div className="text-sm text-gray-500 dark:text-gray-400">Allowed Roles</div>
          <div className="flex flex-wrap gap-2 mt-2">
            {menu.roles && menu.roles.length > 0 ? (
              menu.roles.map((role) => (
                <span
                  key={role.id}
                  className="px-3 py-1 text-sm rounded-full bg-brand-100 text-brand-700 dark:bg-brand-900/20 dark:text-brand-400"
                >
                  {role.name}
                </span>
              ))
            ) : (
              <span className="text-sm text-gray-500">No roles assigned</span>
            )}
          </div>
        </div>

        <div>
          <div className="text-sm text-gray-500 dark:text-gray-400">Status</div>
          <div className="mt-1">
            <span
              className={`px-2 py-1 text-xs rounded-full ${
                menu.is_active
                  ? "bg-success-100 text-success-700 dark:bg-success-900/20 dark:text-success-400"
                  : "bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-400"
              }`}
            >
              {menu.is_active ? "Active" : "Inactive"}
            </span>
          </div>
        </div>

        <div>
          <div className="text-sm text-gray-500 dark:text-gray-400">Created At</div>
          <div className="mt-1 text-gray-800 dark:text-white">
            {new Date(menu.created_at).toLocaleString()}
          </div>
        </div>

        <div>
          <div className="text-sm text-gray-500 dark:text-gray-400">Updated At</div>
          <div className="mt-1 text-gray-800 dark:text-white">
            {new Date(menu.updated_at).toLocaleString()}
          </div>
        </div>

        <div className="flex gap-3 pt-4 border-t border-gray-200 dark:border-gray-800">
          <Link href={`/menus/${menu.id}/edit`}>
            <Button>Edit Menu</Button>
          </Link>
          <Link href="/menus">
            <Button className="bg-gray-500 hover:bg-gray-600">Back to List</Button>
          </Link>
        </div>
      </div>
    </div>
  );
}
