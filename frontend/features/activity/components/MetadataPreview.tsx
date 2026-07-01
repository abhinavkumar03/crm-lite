"use client";

type Props = {
  metadata: string;
};

export default function MetadataPreview({
  metadata,
}: Props) {
  try {
    const decoded = atob(metadata);

    const json = JSON.parse(decoded);

    return (
      <div className="space-y-2">
        {Object.entries(json).map(
          ([key, value]) => (
            <div
              key={key}
              className="flex gap-2 text-sm"
            >
              <span className="font-medium capitalize text-slate-700">
                {key
                .replaceAll("_", " ")
                .replace(
                "duration seconds",
                "Duration"
                )
                .replace(
                "resource type",
                "Resource Type"
                )
                .replace(
                "file name",
                "File Name"
                )
                .replace(
                "preview",
                "Preview"
                )
                .replace(
                "status",
                "Status"
                )
                .replace(
                "direction",
                "Direction"
                )}
                :
              </span>

              <span className="text-slate-500">
                {key ===
                  "duration_seconds"
                  ? `${value} sec`
                  : String(value)}
              </span>
            </div>
          )
        )}
      </div>
    );
  } catch {
    return (
      <p className="text-sm text-slate-500">
        No metadata available.
      </p>
    );
  }
}