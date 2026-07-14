import type { ReactNode } from 'react'
import { Link, usePage } from '@inertiajs/react'
import {
  ExternalLink,
  FolderKanban,
  House,
  ListTodo,
  LogOut,
  Mail,
  Newspaper,
  PanelsTopLeft,
  Tags,
  Users,
  type LucideIcon,
} from 'lucide-react'

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarInset,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarProvider,
  SidebarRail,
  SidebarSeparator,
  SidebarTrigger,
} from '@/components/ui/sidebar'
import { TooltipProvider } from '@/components/ui/tooltip'
import { routes } from '@/routes'

type AdminLayoutProps = {
  children: ReactNode
}

type NavigationItem = {
  label: string
  href: string
  icon: LucideIcon
}

const navigation: NavigationItem[] = [
  { label: 'Home', href: routes.adminDashboardIndex(), icon: House },
  { label: 'Articles', href: routes.adminArticleIndex(), icon: Newspaper },
  { label: 'Newsletters', href: routes.adminNewsletterIndex(), icon: Mail },
  { label: 'Jobs', href: routes.adminJobIndex(), icon: ListTodo },
  { label: 'Projects', href: routes.adminProjectIndex(), icon: FolderKanban },
  { label: 'Subscribers', href: routes.adminSubscriberIndex(), icon: Users },
  { label: 'Tags', href: routes.adminTagIndex(), icon: Tags },
]

function AdminSidebar() {
  const { url } = usePage()
  const pathname = url.split('?')[0]

  function isActive(href: string) {
    return (
      pathname === href ||
      (href !== routes.adminDashboardIndex() && pathname.startsWith(`${href}/`))
    )
  }

  return (
    <Sidebar collapsible="icon">
      <SidebarHeader className="h-16 justify-center border-b border-sidebar-border">
        <div className="flex items-center gap-3 px-3 group-data-[collapsible=icon]:justify-center group-data-[collapsible=icon]:px-0">
          <span className="flex size-8 shrink-0 items-center justify-center bg-sidebar-primary text-sidebar-primary-foreground">
            <PanelsTopLeft className="size-4" />
          </span>
          <div className="min-w-0 leading-tight group-data-[collapsible=icon]:hidden">
            <p className="truncate text-sm font-semibold">Andurel</p>
            <p className="truncate text-xs text-sidebar-foreground/60">Administration</p>
          </div>
        </div>
      </SidebarHeader>

      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupLabel>Navigation</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              {navigation.map((item) => (
                <SidebarMenuItem key={item.href}>
                  <SidebarMenuButton
                    isActive={isActive(item.href)}
                    tooltip={item.label}
                    render={<Link href={item.href} />}
                  >
                    <item.icon />
                    <span>{item.label}</span>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>

      <SidebarSeparator />
      <SidebarFooter>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton tooltip="View site" render={<a href={routes.homePage()} />}>
              <ExternalLink />
              <span>View site</span>
            </SidebarMenuButton>
          </SidebarMenuItem>
          <SidebarMenuItem>
            <SidebarMenuButton
              tooltip="Log out"
              render={<Link href={routes.sessionDestroy()} method="delete" as="button" />}
            >
              <LogOut />
              <span>Log out</span>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  )
}

export default function AdminLayout({ children }: AdminLayoutProps) {
  return (
    <TooltipProvider>
      <SidebarProvider className="h-svh min-h-0 overflow-hidden">
        <AdminSidebar />
        <SidebarInset className="min-h-0 min-w-0 overflow-hidden">
          <header className="flex h-16 shrink-0 items-center border-b border-border px-3">
            <SidebarTrigger />
          </header>
          <div className="min-h-0 flex-1 overflow-auto">{children}</div>
        </SidebarInset>
      </SidebarProvider>
    </TooltipProvider>
  )
}
