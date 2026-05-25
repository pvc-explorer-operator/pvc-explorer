/**
 * Mock snapshot payload for local dev (npm run dev).
 * Covers every phase, mount strategy, access mode, scope, namespace, label combo.
 */

function ts(daysAgo: number): string {
  const d = new Date();
  d.setDate(d.getDate() - daysAgo);
  return d.toISOString();
}

function pvcExplorer(
  name: string,
  namespace: string,
  pvcName: string,
  phase: string,
  mountStrategy: string,
  accessMode: string,
  scope: string | null,
  consumerCount: number,
  labels: Record<string, string>,
  daysAgo: number,
  idleTimeout = '30m',
) {
  const consumers = Array.from({ length: consumerCount }, (_, i) => ({
    podName: `${name}-pod-${i}`,
    ownerKind: i % 2 === 0 ? 'Deployment' : 'StatefulSet',
    ownerName: `${name}-owner-${i}`,
    nodeName: `node-${(i % 3) + 1}`,
    mountReadOnly: accessMode === 'ReadOnlyMany',
  }));

  return {
    metadata: {
      name,
      namespace,
      creationTimestamp: ts(daysAgo),
      labels: {
        ...(scope ? { 'pvcexplorer.io/scope': scope } : {}),
        ...labels,
      },
    },
    spec: {
      pvcName,
      scaling: { idleTimeout },
    },
    status: {
      phase,
      conditions: [
        {
          type: 'Ready',
          status: phase === 'Running' ? 'True' : 'False',
          reason: phase,
          message: `Explorer is ${phase}`,
          lastTransitionTime: ts(daysAgo),
        },
      ],
      mount: {
        strategy: mountStrategy,
        accessMode,
        consumers,
      },
    },
  };
}

function pvcScope(
  name: string,
  phase: string,
  namespaceCount: number,
  explorerCount: number,
  namespaces: string[],
  discoveryMode = 'Label',
  deletionPolicy = 'Retain',
) {
  return {
    metadata: {
      name,
      creationTimestamp: ts(5),
      finalizers: ['pvcexplorer.io/scope-protection'],
    },
    spec: {
      deletionPolicy,
      discovery: { mode: discoveryMode },
      namespaces: { names: namespaces },
      scaling: { idleTimeout: '30m' },
    },
    status: {
      conditions: [
        {
          type: 'Ready',
          status: phase === 'Ready' ? 'True' : 'False',
          reason: phase,
          message: `Scope is ${phase}`,
          lastTransitionTime: ts(5),
        },
      ],
      namespaceCount,
      explorerCount,
    },
  };
}

export const mockScopes = [
  pvcScope('team-alpha',   'Ready',    3, 6, ['team-alpha', 'team-alpha-staging', 'team-alpha-dev']),
  pvcScope('team-beta',    'Ready',    2, 4, ['team-beta', 'team-beta-staging']),
  pvcScope('ml-workloads', 'Ready',    1, 3, ['ml-ns'], 'Label', 'Delete'),
  pvcScope('degraded-env', 'Degraded', 1, 1, ['degraded']),
];

export const mockExplorers = [
  // ── Running ─────────────────────────────────────────────
  pvcExplorer('alpha-data-explorer',    'team-alpha', 'alpha-data-pvc',    'Running',      'Mounted',      'ReadWriteOnce', 'team-alpha',   2, { app: 'analytics', env: 'prod'    }, 10),
  pvcExplorer('alpha-logs-explorer',    'team-alpha', 'alpha-logs-pvc',    'Running',      'Mounted',      'ReadWriteMany', 'team-alpha',   3, { app: 'logging',   env: 'prod'    }, 7),
  pvcExplorer('beta-ml-explorer',       'team-beta',  'beta-ml-pvc',       'Running',      'Mounted',      'ReadWriteOnce', 'team-beta',    1, { app: 'ml',        tier: 'train'  }, 3),
  pvcExplorer('beta-shared-explorer',   'team-beta',  'beta-shared-pvc',   'Running',      'Mounted',      'ReadOnlyMany',  'team-beta',    5, { app: 'shared',    tier: 'serve'  }, 1),

  // ── ScaledToZero ────────────────────────────────────────
  pvcExplorer('alpha-archive-explorer', 'team-alpha', 'alpha-archive-pvc', 'ScaledToZero', 'Unmounted',    'ReadWriteOnce', 'team-alpha',   0, { app: 'archive',   env: 'prod'    }, 30),
  pvcExplorer('beta-cold-explorer',     'team-beta',  'beta-cold-pvc',     'ScaledToZero', 'Unmounted',    'ReadWriteOnce', 'team-beta',    0, { app: 'cold',      tier: 'store'  }, 45),
  pvcExplorer('ml-checkpoint-explorer', 'ml-ns',      'ml-checkpoint-pvc', 'ScaledToZero', 'Unmounted',    'ReadWriteMany', 'ml-workloads', 0, { app: 'ml',        stage: 'ckpt'  }, 14),

  // ── Waking ──────────────────────────────────────────────
  pvcExplorer('alpha-waking-explorer',  'team-alpha', 'alpha-waking-pvc',  'Waking',       'Mounting',     'ReadWriteOnce', 'team-alpha',   0, { app: 'analytics', env: 'staging' }, 5),
  pvcExplorer('ml-waking-explorer',     'ml-ns',      'ml-waking-pvc',     'Waking',       'Mounting',     'ReadWriteOnce', 'ml-workloads', 0, { app: 'ml',        stage: 'eval'  }, 2),

  // ── Pending ─────────────────────────────────────────────
  pvcExplorer('beta-pending-explorer',  'team-beta',  'beta-pending-pvc',  'Pending',      'Pending',      'ReadWriteOnce', 'team-beta',    0, { app: 'ml',        env: 'staging' }, 0),
  pvcExplorer('infra-pending-explorer', 'infra',      'infra-pending-pvc', 'Pending',      'Pending',      'ReadWriteOnce', null,           0, { team: 'infra'                    }, 0),

  // ── Failed ──────────────────────────────────────────────
  pvcExplorer('alpha-failed-explorer',  'team-alpha', 'alpha-failed-pvc',  'Failed',       'Failed',       'ReadWriteOnce', 'team-alpha',   0, { app: 'analytics', env: 'prod'    }, 2),
  pvcExplorer('degraded-explorer',      'degraded',   'degraded-pvc',      'Failed',       'Failed',       'ReadWriteOnce', 'degraded-env', 0, { team: 'ops'                      }, 1),

  // ── InUse (Running + consumers > 0, covered above; add extra) ─
  pvcExplorer('ml-active-explorer',     'ml-ns',      'ml-active-pvc',     'Running',      'Mounted',      'ReadWriteMany', 'ml-workloads', 4, { app: 'ml',        stage: 'train', env: 'prod' }, 1, '1h'),
  pvcExplorer('infra-nfs-explorer',     'infra',      'infra-nfs-pvc',     'Running',      'Mounted',      'ReadOnlyMany',  null,           2, { team: 'infra',    type: 'nfs'    }, 60),

  // ── No scope (covers unscoped filter) ───────────────────
  pvcExplorer('orphan-explorer',        'default',    'orphan-pvc',        'ScaledToZero', 'Unmounted',    'ReadWriteOnce', null,           0, { app: 'orphan'                    }, 90),
  pvcExplorer('solo-running-explorer',  'default',    'solo-pvc',          'Running',      'Mounted',      'ReadWriteOnce', null,           1, { app: 'solo',      env: 'dev'     }, 20),
];
