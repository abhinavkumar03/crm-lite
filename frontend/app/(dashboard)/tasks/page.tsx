"use client";

import { useEffect, useState } from "react";

import {
  createTask,
  deleteTask,
  getTasks,
  updateTask,
} from "@/features/tasks/api";

import {
  getLeads,
} from "@/features/leads/api";

import {
  getContacts,
} from "@/features/contacts/api";

import { Task } from "@/features/tasks/types";
import { Lead } from "@/features/leads/types";
import { Contact } from "@/features/contacts/types";

import TaskTable from "@/features/tasks/components/TaskTable";
import TaskForm from "@/features/tasks/components/TaskForm";

import Modal from "@/components/common/Modal";

import PageHeader from "@/components/common/PageHeader";
import Toolbar from "@/components/common/Toolbar";
import SearchInput from "@/components/common/SearchInput";
import CreateButton from "@/components/common/CreateButton";

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
  }, [page, search]);

  useEffect(() => {
    loadLeads();
    loadContacts();
    }, []);

  async function loadData() {
    await Promise.all([
      loadTasks(),
      loadLeads(),
      loadContacts(),
    ]);
  }

  async function loadTasks() {
  const res = await getTasks({
    page,
    limit: 10,
    search,
  });

  setTasks(res.data.data);
}

 async function loadLeads() {
  const res = await getLeads({
    page: 1,
    limit: 100,
  });

  setLeads(res.data.data);
}

  async function loadContacts() {
  const res = await getContacts({
    page: 1,
    limit: 100,
  });

  setContacts(res.data.data);
}

  async function handleDelete(
    task: Task
  ) {
    const ok =
      window.confirm(
        "Delete this task?"
      );

    if (!ok) return;

    await deleteTask(task.id);

    loadTasks();
  }

  return (
    <div className="space-y-8">
      <PageHeader
        // badge="CRM Module"
        title="Tasks"
        description="Organize work, track progress and connect tasks with leads and contacts."
        action={
          <CreateButton
            text="New Task"
            onClick={() =>
              setOpen(true)
            }
          />
        }
      />

      <Toolbar
        search={
          <SearchInput
            value={search}
            onChange={(value) => {
                setSearch(value);
                setPage(1);
            }}
            placeholder="Search tasks..."
            />
        }
      />

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
        onClose={() =>
          setOpen(false)
        }
      >
        <TaskForm
          submitText="Create Task"
          leads={leads}
          contacts={contacts}
          onSubmit={async (
            values
          ) => {
            await createTask(
              values
            );

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
              title:
                editingTask.title,
              description:
                editingTask.description,
              status:
                editingTask.status,
              due_date:
                editingTask.due_date,
              lead_id:
                editingTask.lead_id,
              contact_id:
                editingTask.contact_id,
            }}
            submitText="Update Task"
            leads={leads}
            contacts={contacts}
            onSubmit={async (
              values
            ) => {
              await updateTask(
                editingTask.id,
                values
              );

              setEditingTask(
                null
              );

              loadTasks();
            }}
          />
        )}
      </Modal>
    </div>
  );
}