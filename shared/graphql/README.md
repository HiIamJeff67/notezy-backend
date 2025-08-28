### GraphQL for Notezy Backend, Frontend, and Software

- Folders which have used in backend:

  - schemas/
  - scalars/

- Folders which have used in web frontend:

  - schemas/
  - fragments/
  - queries/

- Note:
  - This folder will be built to submodule one day.
  - The schemas(`*.graphql` files) are the same in all the environments, and should be maintain or synchronize them periodically, but the yaml file to generate the adapted codes are different, for example: the backend use `gqlgen.yml` to generate codes, and the web frontend use `codegen.yml` otherwise.
  - This submodule MUST be a private property of Notezy, this means don't release this submodule and don't make it public in github.
