<script lang="ts">
  import type AttendeeType from "./attendee";
  import Attendee from "./Attendee.svelte";

  async function listAttendees(): Promise<AttendeeType[]> {
    return fetch("/attendees").then((response) => {
      if (!response.ok) {
        throw new Error(response.statusText);
      }
      return response.json() as Promise<AttendeeType[]>;
    });
  }

  let attendeesPromise: Promise<AttendeeType[]> = listAttendees();
</script>

<main>
  <h1>Attendance List</h1>
  {#await attendeesPromise}
    <p>...waiting</p>
  {:then attendees}
    <div class="attendees-container">
      {#each attendees as attendee}
        <Attendee {attendee} />
      {/each}
    </div>
  {:catch error}
    <p style="color: red">{error}</p>
  {/await}
</main>

<style>
  .attendees-container {
    display: flex;
    flex-wrap: wrap;
  }
</style>
