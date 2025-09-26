-- Drop tables in reverse order of creation
DROP TABLE IF EXISTS public.waitlists;
DROP TABLE IF EXISTS public.users;

-- Drop enum types
DROP TYPE IF EXISTS public.auth_provider;