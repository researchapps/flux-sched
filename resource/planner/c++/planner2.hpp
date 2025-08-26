/*****************************************************************************\
 * Copyright 2025 Lawrence Livermore National Security, LLC
 * (c.f. AUTHORS, NOTICE.LLNS, LICENSE)
 *
 * This file is part of the Flux resource manager framework.
 * For details, see https://github.com/flux-framework.
 *
 * SPDX-License-Identifier: LGPL-3.0
\*****************************************************************************/

#ifndef PLANNER_HPP
#define PLANNER_HPP

#include <string>
#include <map>
#include <unordered_map>
#include <boost/icl/interval_map.hpp>
#include <boost/icl/interval.hpp>

// Use boost's interval map for our resource plan.
// The key is a discrete interval of time.
// The value is the number of free resources at that time.
using resource_plan_t = boost::icl::interval_map<uint64_t, uint64_t>;
using time_interval = boost::icl::discrete_interval<uint64_t>;

/*! A struct to store information about a planned resource span.
 */
struct span_t {
    uint64_t start = 0;           /* start time of the span */
    uint64_t end = 0;             /* end time of the span (exclusive) */
    int64_t span_id = 0;          /* unique span id */
    uint64_t res_occupied = 0;    /* required resource quantity */
};

/*! The planner class, rewritten to use Boost.ICL.
 */
class planner2 {
public:
    planner2 () = default;
    planner2 (const uint64_t total_resources,
             const std::string &resource_type,
             const uint64_t plan_start,
             const uint64_t plan_end);
    ~planner2 () = default;

    // Public methods for managing the plan.
    int64_t add_span (uint64_t start_time, uint64_t duration, uint64_t request);
    int64_t remove_span (int64_t span_id);

    // Public methods for querying the plan.
    bool avail_during (uint64_t at, uint64_t duration, uint64_t request) const;
    int64_t avail_time_first (uint64_t at, uint64_t duration, uint64_t request) const;
    int64_t avail_resources_during (uint64_t at, uint64_t duration) const;
    int64_t avail_resources_at (uint64_t at) const;

private:
    // The Boost.ICL interval map to store the resource plan.
    resource_plan_t m_plan;

    // Metadata about the plan.
    uint64_t m_total_resources = 0;
    std::string m_resource_type = "";
    uint64_t m_plan_start = 0;
    uint64_t m_plan_end = 0;
    uint64_t m_span_counter = 0;

    // A map to look up spans by their ID.
    std::map<int64_t, span_t> m_span_lookup;
};

#endif // PLANNER_HPP
