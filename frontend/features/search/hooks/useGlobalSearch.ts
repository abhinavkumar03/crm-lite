"use client";

import { useEffect, useState } from "react";

import { searchCRM } from "../api";
import { SearchResponse } from "../types";

const empty: SearchResponse = { results: [] };

export default function useGlobalSearch(query: string) {
  const [loading, setLoading] = useState(false);
  const [results, setResults] = useState<SearchResponse>(empty);

  useEffect(() => {
    if (query.trim().length < 2) {
      setResults(empty);
      return;
    }

    const timeout = setTimeout(async () => {
      try {
        setLoading(true);
        const data = await searchCRM(query);
        setResults(data ?? empty);
      } catch {
        setResults(empty);
      } finally {
        setLoading(false);
      }
    }, 350);

    return () => clearTimeout(timeout);
  }, [query]);

  return { loading, results };
}
