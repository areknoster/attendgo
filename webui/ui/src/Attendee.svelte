<script lang="ts">
  import type Attendee from "./attendee";
  export let attendee: Attendee;
</script>

<figure>
  <img
    src="/photos/{attendee.PhotoId}.jpg"
    alt="attendee's {attendee.Id} face"
  />
  <figcaption>{attendee.Id}</figcaption>
</figure>

<style>
  figure {
    max-width: 45vw;
  }

  figure img {
    max-width: 90%;
    padding: auto;
  }

  figure figcaption {
    text-align: center;
    size: 10em;
  }

</style>
