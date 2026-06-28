"use client";

import { useEffect, useState } from "react";

import Modal from "@/components/common/Modal";

import {
    createTask,
    deleteTask,
    getTasks,
    updateTask,
} from "@/features/tasks/api";

import { Task } from "@/features/tasks/types";

import TaskTable from "@/features/tasks/components/TaskTable";

import TaskForm from "@/features/tasks/components/TaskForm";

import { getLeads } from "@/features/leads/api";
import { Lead } from "@/features/leads/types";

import { getContacts } from "@/features/contacts/api";
import { Contact } from "@/features/contacts/types";

export default function TasksPage() {

    const [tasks, setTasks] =
        useState<Task[]>([]);

    const [leads, setLeads] =
        useState<Lead[]>([]);

    const [contacts, setContacts] =
        useState<Contact[]>([]);

    const [page, setPage] =
        useState(1);

    const [search, setSearch] =
        useState("");

    const [open, setOpen] =
        useState(false);

    const [editingTask, setEditingTask] =
        useState<Task | null>(null);

    useEffect(() => {

        loadData();

    }, [page]);

    async function loadData() {

        await Promise.all([
            loadTasks(),
            loadLeads(),
            loadContacts(),
        ]);

    }

    async function loadTasks() {

        const res = await getTasks(
            page,
            search,
        );

        setTasks(res.data);

    }

    async function loadLeads() {

        const res =
            await getLeads(1, "");

        setLeads(res.data);

    }

    async function loadContacts() {

        const res =
            await getContacts(1, "");

        setContacts(res.data);

    }

    async function handleDelete(
        task: Task,
    ) {

        const ok = window.confirm(
            "Delete this task?"
        );

        if (!ok) return;

        await deleteTask(task.id);

        await loadTasks();

    }

    return (

        <div className="space-y-6">

            <div className="flex items-center justify-between">

                <h1 className="text-3xl font-bold">

                    Tasks

                </h1>

                <button
                    onClick={() => setOpen(true)}
                    className="rounded bg-blue-600 px-4 py-2 text-white"
                >

                    Create Task

                </button>

            </div>

            <div className="flex gap-2">

                <input
                    value={search}
                    onChange={(e) =>
                        setSearch(e.target.value)
                    }
                    placeholder="Search..."
                    className="flex-1 rounded border p-3"
                />

                <button
                    onClick={() => {
                        setPage(1);
                        loadTasks();
                    }}
                    className="rounded bg-blue-600 px-4 py-2 text-white"
                >

                    Search

                </button>

            </div>

            <TaskTable
                tasks={tasks}
                page={page}
                setPage={setPage}
                onEdit={setEditingTask}
                onDelete={handleDelete}
            />

            <Modal
                open={open}
                title="Create Task"
                onClose={() => setOpen(false)}
            >

                <TaskForm
                    submitText="Create Task"
                    leads={leads}
                    contacts={contacts}
                    onSubmit={async (values) => {

                        await createTask(values);

                        setOpen(false);

                        loadTasks();

                    }}
                />

            </Modal>

            <Modal
                open={!!editingTask}
                title="Edit Task"
                onClose={() =>
                    setEditingTask(null)
                }
            >

                {editingTask && (

                    <TaskForm
                        initialValues={{
                            title: editingTask.title,
                            description: editingTask.description,
                            status: editingTask.status,
                            due_date: editingTask.due_date,
                            lead_id: editingTask.lead_id,
                            contact_id: editingTask.contact_id,
                        }}
                        submitText="Update Task"
                        leads={leads}
                        contacts={contacts}
                        onSubmit={async (values) => {

                            await updateTask(
                                editingTask.id,
                                values,
                            );

                            setEditingTask(null);

                            loadTasks();

                        }}
                    />

                )}

            </Modal>

        </div>

    );

}