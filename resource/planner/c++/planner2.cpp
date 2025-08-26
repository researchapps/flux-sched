/*****************************************************************************\
 * Copyright 2025 Lawrence Livermore National Security, LLC
 * (c.f. AUTHORS, NOTICE.LLNS, LICENSE)
 *
 * This file is part of the Flux resource manager framework.
 * For details, see https://github.com/flux-framework.
 *
 * SPDX-License-Identifier: LGPL-3.0
\*****************************************************************************/

#include <iostream>
#include <errno.h>

#include "planner2.hpp"

/****************************************************************************
 * *
 * Public Planner Methods                               *
 * *
 ****************************************************************************/

planner2::planner2 (const uint64_t total_resources,
                  const std::string &resource_type,
                  const uint64_t plan_start,
                  const uint64_t plan_end)
{
    m_total_resources = total_resources;
    m_resource_type = resource_type;
    m_plan_start = plan_start;
    m_plan_end = plan_end;
    
    // Initialize the entire plan with all resources free.
    m_plan.add(std::make_pair(time_interval::right_open(m_plan_start, m_plan_end), m_total_resources));
}

bool planner2::avail_during (uint64_t at, uint64_t duration, uint64_t request) const
{
    if ((at + duration) > m_plan_end || request > m_total_resources) {
        errno = EINVAL;
        return false;
    }

    // Get the submap for the requested interval
    time_interval occupied_interval = time_interval::right_open(at, at + duration);
    resource_plan_t submap = m_plan & occupied_interval;

    // Iterate through the intervals in the submap
    for (auto const& pair : submap) {
        if (pair.second < request) {
            return false;
        }
    }
    
    return true;
}

int64_t planner2::add_span (uint64_t start_time, uint64_t duration, uint64_t request)
{
    if (start_time < m_plan_start || duration == 0 ||
        (start_time + duration) > m_plan_end || request > m_total_resources) {
        errno = EINVAL;
        return -1;
    }
    
    // Check for availability before adding
    if (!avail_during(start_time, duration, request)) {
        errno = EPERM;
        return -1;
    }

    try {
        ++m_span_counter;
        if (m_span_lookup.find(m_span_counter) != m_span_lookup.end()) {
            errno = EEXIST;
            return -1;
        }

        span_t new_span;
        new_span.span_id = m_span_counter;
        new_span.start = start_time;
        new_span.end = start_time + duration;
        new_span.res_occupied = request;
        
        // Store the span metadata
        m_span_lookup.insert({new_span.span_id, new_span});

        // Use Boost.ICL to subtract the requested resources over the interval
        time_interval occupied_interval = time_interval::right_open(start_time, start_time + duration);
        m_plan -= std::make_pair(occupied_interval, request);
    } catch (std::bad_alloc& e) {
        errno = ENOMEM;
        return -1;
    }
    
    return m_span_counter;
}

int64_t planner2::remove_span (int64_t span_id)
{
    auto it = m_span_lookup.find(span_id);
    if (it == m_span_lookup.end()) {
        errno = ENOENT;
        return -1;
    }
    
    span_t span = it->second;
    
    // Use Boost.ICL to add the resources back to the interval
    time_interval freed_interval = time_interval::right_open(span.start, span.end);
    m_plan += std::make_pair(freed_interval, span.res_occupied);
    
    // Remove the span from our lookup map
    m_span_lookup.erase(it);

    return 0;
}

int64_t planner2::avail_time_first (uint64_t at, uint64_t duration, uint64_t request) const
{
    if (duration == 0 || request > m_total_resources || at >= m_plan_end) {
        errno = EINVAL;
        return -1;
    }

    // First, check if the requested time 'at' is itself a valid starting point for the duration.
    if (avail_during(at, duration, request)) {
        return at;
    }

    // If not, perform a logarithmic search for the first available slot after 'at'.
    // The upper_bound call finds the first interval that starts strictly after 'at'.
    auto it = m_plan.upper_bound(time_interval::right_open(at, at));

    // Iterate from the first interval found by the logarithmic search.
    for (; it != m_plan.end(); ++it) {
        uint64_t start_of_slot = boost::icl::lower(it->first);
        // Check if this interval and the following ones can accommodate the request.
        if (avail_during(start_of_slot, duration, request)) {
            return start_of_slot;
        }
    }
    
    errno = ENOENT;
    return -1;
}

int64_t planner2::avail_resources_during (uint64_t at, uint64_t duration) const
{
    if ((at + duration) > m_plan_end) {
        errno = EINVAL;
        return -1;
    }
    
    uint64_t min_free = m_total_resources;

    // Use Boost.ICL to create a sub-map for the requested interval.
    time_interval query_interval = time_interval::right_open(at, at + duration);
    resource_plan_t submap = m_plan & query_interval;

    // Iterate through the resulting intervals and find the minimum free count
    for (auto const& pair : submap) {
        if (pair.second < min_free) {
            min_free = pair.second;
        }
    }

    return min_free;
}

int64_t planner2::avail_resources_at (uint64_t at) const
{
    if (at > m_plan_end || at < m_plan_start) {
        errno = EINVAL;
        return -1;
    }

    // Find the value at a specific point in time.
    auto it = m_plan.find(at);
    if (it != m_plan.end()) {
        return it->second;
    }

    // If the point is not found in any interval, return -1 with an error.
    errno = ENOENT;
    return -1;
}
