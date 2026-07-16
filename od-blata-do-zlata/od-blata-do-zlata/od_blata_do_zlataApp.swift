//
//  od_blata_do_zlataApp.swift
//  od-blata-do-zlata
//
//  Created by Antonio Obradovic on 16.07.2026..
//

import SwiftUI
import SwiftData

@main
struct od_blata_do_zlataApp: App {
    var sharedModelContainer: ModelContainer = {
        let schema = Schema([
            Item.self,
        ])
        let modelConfiguration = ModelConfiguration(schema: schema, isStoredInMemoryOnly: false)

        do {
            return try ModelContainer(for: schema, configurations: [modelConfiguration])
        } catch {
            fatalError("Could not create ModelContainer: \(error)")
        }
    }()

    var body: some Scene {
        WindowGroup {
            ContentView()
        }
        .modelContainer(sharedModelContainer)
    }
}
