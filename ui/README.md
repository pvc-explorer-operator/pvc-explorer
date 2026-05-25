# UI Development

## Run UI without building images

You can test the UI locally with Vite. No Docker image build is required.

```bash
cd ui
npm install
npm run dev
```

## Test login screen locally

By default, development mode auto-authenticates as `devuser` for faster iteration.
To test the real login flow and keep the login page active:

```bash
cd ui
VITE_DEV_AUTH_BYPASS=false npm run dev
```

Then run the backend in another terminal (also without Docker image build):

```bash
make run
```

Open the UI and navigate to `/login`.
